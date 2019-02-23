# from datetime import datetime

from flask import request
from flask_restful import Resource, abort
from flask_restful.reqparse import RequestParser

github = None


def _get_user_email():
    """Get the user object or create it based on the token in the session

    If there is no access token abort with 401 message
    """
    if 'Access-Token' not in request.headers:
        abort(401, message='Access Denied!')

    token = request.headers['Access-Token']
    user_data = github.get('user', token=dict(access_token=token)).data
    email = user_data['email']
    # name = user_data['name']

    return email


class Link(Resource):
    def get(self):
        """Get all links

        If user doesn't exist create it (with no goals)
        """
        user = _get_user_email()

        result = {}
        return result

    def post(self):
        user = _get_user_email()
        parser = RequestParser()
        parser.add_argument('url', type=str, required=True)
        parser.add_argument('title', type=str, required=True)
        parser.add_argument('description', type=str, required=False)
        args = parser.parse_args()
