#!/usr/bin/env python3

from http.server import BaseHTTPRequestHandler, HTTPServer
from urllib.parse import urlparse
import json
from sense_hat import SenseHat

class RequestHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        # SenseHAT stuff
        sense = SenseHat()
        t = sense.get_temperature()
        p = sense.get_pressure()
        h = sense.get_humidity()
        t = round(t)
        p = round(p)
        h = round(h)
        f_t = round(9.0/5.0 * t + 32, 2)

        self.send_response(200)
        self.end_headers()
        self.wfile.write(json.dumps({
            'temperature': f_t,
            'pressure': p,
            'humidity': h
        }).encode())
        return

if __name__ == '__main__':
    server = HTTPServer(('0.0.0.0', 8000), RequestHandler)
    print('Starting server at http://clock.lan:8000')
    server.serve_forever()
