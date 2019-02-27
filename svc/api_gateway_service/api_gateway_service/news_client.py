import grpc
import news_pb2
import news_pb2_grpc


class NewsClient(object):
    """
    Client for accessing the gRPC functionality
    """

    def __init__(self, host, port):
        # instantiate a communication channel
        self.channel = grpc.insecure_channel(
            '{}:{}'.format(host, port))

        # bind the client to the server channel
        self.stub = news_pb2_grpc.NewsStub(self.channel)

    def get_news(self, username, startToken=None):
        """
        Client function to call the rpc for GetDigest
        """
        req = news_pb2.GetNewsRequest(username=username, startToken=startToken)
        return self.stub.GetNews(req)
