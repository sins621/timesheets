#!/usr/bin/env python3
"""
Test HTTP server that serves up API information.
Enhanced with Rich for beautiful CLI output.
"""

import json
import http.server
import socketserver
from urllib.parse import urlparse, parse_qs
from datetime import datetime
import uuid
import sys
import os

if os.path.exists('.venv'):
    sys.path.insert(0, os.path.join('.venv', 'lib', 'python3.13', 'site-packages'))

try:
    from rich.console import Console
    from rich.table import Table
    from rich.panel import Panel
    from rich.text import Text
    from rich.live import Live
    from rich.layout import Layout
    from rich.align import Align
    from rich.columns import Columns
    from rich.progress import Progress, SpinnerColumn, TextColumn
    from rich.logging import RichHandler
    from rich.traceback import install
    import logging
    
    install(show_locals=True)
    logging.basicConfig(
        level=logging.INFO,
        format="%(message)s",
        datefmt="[%X]",
        handlers=[RichHandler(rich_tracebacks=True)]
    )
    
    RICH_AVAILABLE = True
except ImportError:
    RICH_AVAILABLE = False
    import logging
    logging.basicConfig(level=logging.INFO)

console = Console() if RICH_AVAILABLE else None

class APIHandler(http.server.BaseHTTPRequestHandler):
    valid_tokens = set()
    request_count = 0
    
    @classmethod
    def get_stats(cls):
        """Get server statistics"""
        return {
            "active_tokens": len(cls.valid_tokens),
            "total_requests": cls.request_count
        }
    
    def check_bearer_auth(self):
        """Check if request has valid bearer token"""
        auth_header = self.headers.get('Authorization', '')
        if not auth_header.startswith('Bearer '):
            return False
        
        token = auth_header[7:]
        return token in self.valid_tokens
    
    def require_auth(self):
        """Send 401 if no valid bearer token"""
        if not self.check_bearer_auth():
            auth_header = self.headers.get('Authorization', 'None')
            if RICH_AVAILABLE and console:
                if auth_header == 'None':
                    console.print(f"[red]🚫 Auth failed:[/red] [dim]No Authorization header[/dim]")
                else:
                    token_preview = auth_header[7:15] + "..." if len(auth_header) > 15 else auth_header[7:]
                    console.print(f"[red]🚫 Auth failed:[/red] [dim]Invalid token: {token_preview}[/dim]")
            
            self.send_response(401)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            error_response = {
                "error": "Unauthorized",
                "message": "Valid bearer token required"
            }
            self.wfile.write(json.dumps(error_response, indent=2).encode())
            return False
        return True
    
    def do_GET(self):
        """Handle GET requests"""
        parsed_path = urlparse(self.path)
        path = parsed_path.path
        
        if path == '/':
            self.serve_api_docs()
        elif path == '/api/client/list':
            if self.require_auth():
                self.list_customers()
        elif path.startswith('/api/project/client/'):
            if self.require_auth():
                client_id = path.split('/')[-1]
                self.list_projects_by_customer(client_id)
        else:
            self.send_error(404, "Not Found")
    
    def do_POST(self):
        """Handle POST requests"""
        parsed_path = urlparse(self.path)
        path = parsed_path.path
        
        if path == '/api/account/Authorise':
            self.handle_login()
        elif path == '/api/entry/create':
            if self.require_auth():
                self.handle_make_entry()
        else:
            self.send_error(404, "Not Found")
    
    def serve_api_docs(self):
        """Serve the main API documentation page"""
        html_content = """
<!DOCTYPE html>
<html>
<head>
    <title>Warp Development API Test Server</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #1e1e1e; color: #d4d4d4; }
        .endpoint { margin: 20px 0; padding: 15px; background: #2d2d30; border-radius: 5px; }
        .method { color: #4ec9b0; font-weight: bold; }
        .url { color: #9cdcfe; }
        .note { color: #608b4e; font-style: italic; }
        pre { background: #0d1117; padding: 10px; border-radius: 3px; overflow-x: auto; }
        h1 { color: #569cd6; }
        h2 { color: #4ec9b0; }
    </style>
</head>
<body>
    <h1>Warp Development API Test Server</h1>
    <p class="note">(I'm not sure if these APIs work, it hasn't been touched in years)</p>
    
    <div class="endpoint">
        <h2>Login</h2>
        <p><span class="method">POST</span> <span class="url">https://office.warpdevelopment.com/api/account/Authorise</span></p>
        <pre>{
  "Email": "email@email.com",
  "Password": "Password123"
}</pre>
        <p>The following calls probably need cookie auth based on the above.</p>
    </div>
    
    <div class="endpoint">
        <h2>List Customers</h2>
        <p><span class="method">GET</span> <span class="url">https://office.warpdevelopment.com/api/client/list</span></p>
    </div>
    
    <div class="endpoint">
        <h2>List Projects By Customer</h2>
        <p><span class="method">GET</span> <span class="url">https://office.warpdevelopment.com/api/project/client/{clientId}</span></p>
    </div>
    
    <div class="endpoint">
        <h2>Make Entry</h2>
        <p><span class="method">POST</span> <span class="url">https://office.warpdevelopment.com/api/entry/create</span></p>
        <pre>{
  "Comments": "this is a entry",
  "EntryDate": "2019-09-01",
  "Time": 8,
  "Overtime": 0, //1==true
  "Person": {"PersonId": 123},
  "Task": {"TaskId": 123},
  "CostCodeId": 1,
}</pre>
    </div>
    
    <hr>
    <h2>Test Endpoints (Local)</h2>
    <p>This test server provides mock responses for the above endpoints:</p>
    <p><strong>Authentication:</strong> Use POST /api/account/Authorise to get a bearer token, then include it in the Authorization header as "Bearer {token}" for protected endpoints.</p>
    <ul>
        <li>POST /api/account/Authorise - Mock login (accepts any credentials, returns bearer token)</li>
        <li>GET /api/client/list - Mock customer list (requires bearer token)</li>
        <li>GET /api/project/client/123 - Mock project list for customer 123 (requires bearer token)</li>
        <li>POST /api/entry/create - Mock entry creation (requires bearer token)</li>
    </ul>
</body>
</html>
        """
        
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()
        self.wfile.write(html_content.encode())
    
    def handle_login(self):
        """Handle login POST request"""
        content_length = int(self.headers['Content-Length'])
        post_data = self.rfile.read(content_length)
        
        try:
            login_data = json.loads(post_data.decode('utf-8'))
            email = login_data.get('Email', '')
            password = login_data.get('Password', '')
            
            bearer_token = str(uuid.uuid4())
            
            self.valid_tokens.add(bearer_token)
            
            if RICH_AVAILABLE and console:
                console.print(f"[green]🔐 New login:[/green] [cyan]{email}[/cyan] [dim]→ token: {bearer_token[:8]}...[/dim]")
            response = {
                "success": True,
                "message": "Login successful",
                "token": bearer_token,
                "token_type": "Bearer",
                "user": {
                    "email": email,
                    "userId": 123
                }
            }
            
            self.send_response(200)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(json.dumps(response, indent=2).encode())
            
        except json.JSONDecodeError:
            self.send_error(400, "Invalid JSON")
    
    def list_customers(self):
        """Handle list customers GET request"""
        mock_customers = [
            {
                "clientId": 123,
                "name": "Acme Corporation",
                "email": "contact@acme.com",
                "phone": "+1-555-0123",
                "active": True
            },
            {
                "clientId": 456,
                "name": "Tech Solutions Ltd",
                "email": "info@techsolutions.com",
                "phone": "+1-555-0456",
                "active": True
            },
            {
                "clientId": 789,
                "name": "Global Industries",
                "email": "hello@globalind.com",
                "phone": "+1-555-0789",
                "active": False
            }
        ]
        
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps(mock_customers, indent=2).encode())
    
    def list_projects_by_customer(self, client_id):
        """Handle list projects by customer GET request"""
        mock_projects = [
            {
                "projectId": 1001,
                "clientId": int(client_id) if client_id.isdigit() else 123,
                "name": "Website Redesign",
                "description": "Complete overhaul of company website",
                "status": "active",
                "startDate": "2024-01-15",
                "estimatedHours": 120
            },
            {
                "projectId": 1002,
                "clientId": int(client_id) if client_id.isdigit() else 123,
                "name": "Mobile App Development",
                "description": "Native iOS and Android application",
                "status": "planning",
                "startDate": "2024-03-01",
                "estimatedHours": 300
            }
        ]
        
        self.send_response(200)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps(mock_projects, indent=2).encode())
    
    def handle_make_entry(self):
        """Handle make entry POST request"""
        content_length = int(self.headers['Content-Length'])
        post_data = self.rfile.read(content_length)
        
        try:
            entry_data = json.loads(post_data.decode('utf-8'))
            

            response = {
                "success": True,
                "entryId": 9876,
                "message": "Entry created successfully",
                "entry": {
                    "entryId": 9876,
                    "comments": entry_data.get("Comments", ""),
                    "entryDate": entry_data.get("EntryDate", datetime.now().strftime("%Y-%m-%d")),
                    "time": entry_data.get("Time", 0),
                    "overtime": entry_data.get("Overtime", 0),
                    "personId": entry_data.get("Person", {}).get("PersonId", 123),
                    "taskId": entry_data.get("Task", {}).get("TaskId", 123),
                    "costCodeId": entry_data.get("CostCodeId", 1),
                    "createdAt": datetime.now().isoformat()
                }
            }
            
            self.send_response(201)
            self.send_header('Content-type', 'application/json')
            self.end_headers()
            self.wfile.write(json.dumps(response, indent=2).encode())
            
        except json.JSONDecodeError:
            self.send_error(400, "Invalid JSON")
    
    def log_message(self, format, *args):
        """Custom log message format with Rich styling"""
        self.__class__.request_count += 1
        
        if RICH_AVAILABLE and console:
            message = format % args
            parts = message.split(' ')
            
            if len(parts) >= 3:
                method = parts[0].strip('"')
                path = parts[1]
                status = parts[2]
                

                method_colors = {
                    'GET': 'green',
                    'POST': 'blue', 
                    'PUT': 'yellow',
                    'DELETE': 'red',
                    'PATCH': 'magenta'
                }
                

                if status.startswith('2'):
                    status_color = 'green'
                elif status.startswith('4'):
                    status_color = 'yellow'
                elif status.startswith('5'):
                    status_color = 'red'
                else:
                    status_color = 'white'
                
                method_color = method_colors.get(method, 'white')
                
                console.print(
                    f"[dim]{datetime.now().strftime('%H:%M:%S')}[/dim] "
                    f"[{method_color}]{method:>6}[/{method_color}] "
                    f"[cyan]{path}[/cyan] "
                    f"[{status_color}]{status}[/{status_color}] "
                    f"[dim]{self.address_string()}[/dim]"
                )
            else:
                console.print(f"[dim]{datetime.now().strftime('%H:%M:%S')}[/dim] {message}")
        else:
            print(f"[{datetime.now().strftime('%Y-%m-%d %H:%M:%S')}] {format % args}")

def create_startup_panel():
    """Create a beautiful startup panel"""
    if not RICH_AVAILABLE:
        return None
        
    startup_text = Text()
    startup_text.append("🚀 Warp Development API Test Server\n", style="bold blue")
    startup_text.append("Enhanced with Rich for beautiful output\n\n", style="dim")
    
    startup_text.append("📋 Available Endpoints:\n", style="bold green")
    startup_text.append("  • POST /api/account/Authorise - Login & get bearer token\n", style="cyan")
    startup_text.append("  • GET  /api/client/list - List customers (auth required)\n", style="cyan")
    startup_text.append("  • GET  /api/project/client/{id} - List projects (auth required)\n", style="cyan")
    startup_text.append("  • POST /api/entry/create - Create time entry (auth required)\n", style="cyan")
    startup_text.append("  • GET  / - API documentation page\n\n", style="cyan")
    
    startup_text.append("🔗 Quick Links:\n", style="bold green")
    startup_text.append("  • Documentation: http://localhost:8000\n", style="yellow")
    startup_text.append("  • Postman Collection: Warp-API-Test.postman_collection.json\n", style="yellow")
    startup_text.append("\n💡 Press Ctrl+C to stop the server", style="dim")
    
    return Panel(
        startup_text,
        title="[bold white]🌟 Server Starting[/bold white]",
        border_style="blue",
        padding=(1, 2)
    )

def create_stats_table():
    """Create a live stats table"""
    if not RICH_AVAILABLE:
        return None
        
    stats = APIHandler.get_stats()
    
    table = Table(show_header=True, header_style="bold magenta", border_style="blue")
    table.add_column("Metric", style="cyan", no_wrap=True)
    table.add_column("Value", style="green", justify="right")
    
    table.add_row("🔑 Active Tokens", str(stats["active_tokens"]))
    table.add_row("📊 Total Requests", str(stats["total_requests"]))
    table.add_row("⏰ Server Uptime", datetime.now().strftime("%H:%M:%S"))
    
    return Panel(
        table,
        title="[bold white]📈 Server Stats[/bold white]",
        border_style="green",
        padding=(0, 1)
    )

def main():
    PORT = 8000
    
    if RICH_AVAILABLE and console:
        console.clear()
        startup_panel = create_startup_panel()
        if startup_panel:
            console.print(startup_panel)
            console.print()
        
        console.print(f"[bold green]✅ Server running on[/bold green] [bold blue]http://localhost:{PORT}[/bold blue]")
        console.print("[dim]Logs will appear below...[/dim]\n")
    else:
        print(f"Starting test HTTP server on port {PORT}")
        print(f"Visit http://localhost:{PORT} to see the API documentation")
        print("Press Ctrl+C to stop the server")
    
    try:
        with socketserver.TCPServer(("", PORT), APIHandler) as httpd:
            httpd.serve_forever()
    except KeyboardInterrupt:
        if RICH_AVAILABLE and console:
            console.print("\n[bold red]🛑 Server stopped.[/bold red]")
            
            final_stats = APIHandler.get_stats()
            console.print(f"[dim]Final stats: {final_stats['total_requests']} requests served, {final_stats['active_tokens']} tokens issued[/dim]")
        else:
            print("\nServer stopped.")

if __name__ == "__main__":
    main()
