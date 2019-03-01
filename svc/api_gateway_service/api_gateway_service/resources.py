# from datetime import datetime
import os

import requests
from flask import request
from flask_restful import Resource, abort
from flask_restful.reqparse import RequestParser

github = None


def _get_user():
    """Get the user object or create it based on the token in the session

    If there is no access token abort with 401 message
    """
    if 'Access-Token' not in request.headers:
        abort(401, message='Access Denied!')

    token = request.headers['Access-Token']
    user_data = github.get('user', token=dict(access_token=token)).data
    if 'email' not in user_data:
        abort(401, message='Access Denied!')

    email = user_data['email']
    name = user_data['name']

    return name, email


class Link(Resource):
    host = os.environ.get('LINK_MANAGER_SERVICE_HOST', 'localhost')
    port = os.environ.get('LINK_MANAGER_SERVICE_PORT', '8080')
    base_url = 'http://{}:{}/links'.format(host, port)

    def get(self):
        """Get all links

        If user doesn't exist create it (with no goals)
        """
        username, email = _get_user()
        parser = RequestParser()
        parser.add_argument('url_regex', type=str, required=False)
        parser.add_argument('title_regex', type=str, required=False)
        parser.add_argument('description_regex', type=str, required=False)
        parser.add_argument('tag', type=str, required=False)
        parser.add_argument('start_token', type=str, required=False)
        args = parser.parse_args()
        args.update(username=username)
        r = requests.get(self.base_url, params=args)

        if not r.ok:
            abort(r.status_code, message=r.content)

        return r.json()

    def post(self):
        username, email = _get_user()
        parser = RequestParser()
        parser.add_argument('url', type=str, required=True)
        parser.add_argument('title', type=str, required=True)
        parser.add_argument('description', type=str, required=False)
        parser.add_argument('tags', type=str, required=False)
        parser.add_argument('start_token', type=str, required=False)
        args = parser.parse_args()
        args.update(username=username)
        r = requests.post(self.base_url, json=args)
        if not r.ok:
            abort(r.status_code, message=r.content)

        return r.json()

    def put(self):
        username, email = _get_user()
        parser = RequestParser()
        parser.add_argument('url', type=str, required=True)
        parser.add_argument('title', type=str, required=True)
        parser.add_argument('description', type=str, required=False)
        parser.add_argument('add_tags', type=str, required=False)
        parser.add_argument('remove_tags', type=str, required=False)
        args = parser.parse_args()
        args.update(username=username)
        r = requests.put(self.base_url, json=args)
        if not r.ok:
            abort(r.status_code, message=r.content)

        return r.json()

    def delete(self):
        username, email = _get_user()
        parser = RequestParser()
        parser.add_argument('url', type=str, required=True)
        args = parser.parse_args()
        args.update(username=username)
        r = requests.delete(self.base_url, **args)
        if not r.ok:
            abort(r.status_code, message=r.content)

        return r.json()
