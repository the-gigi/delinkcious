import os
from api_gateway_service.api import app


def main():
    port = int(os.environ.get('PORT', 5000))
    login_url = 'http://localhost:{}/login'.format(port)
    print('If you run locally, browse to', login_url)
    host = '0.0.0.0'
    app.run(host=host, port=port)


if __name__ == "__main__":
    main()
