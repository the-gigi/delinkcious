import os
from news_client import NewsClient

if __name__ == '__main__':
    host = os.environ.get('NEWS_MANAGER_SERVICE_HOST', 'localhost')
    port = int(os.environ.get('NEWS_MANAGER_SERVICE_PORT', '6060'))
    cli = NewsClient(host, port)
    resp = cli.get_news('gigi')
    print(resp)