import os
from api_gateway_service.api import app


def main():
    port = int(os.environ.get('PORT', 6000))
    print(f'If you run locally, browse to http://localhost:{port}')
    host = '0.0.0.0'
    app.run(host=host, port=port)


if __name__ == "__main__":
    main()
