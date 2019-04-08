# from datetime import datetime
import os

import requests
from flask import request
from flask_restful import Resource, abort
from flask_restful.reqparse import RequestParser
from urllib.parse import urlencode

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
        url = '{}?{}'.format(self.base_url, urlencode(args))
        r = requests.delete(url)
        if not r.ok:
            abort(r.status_code, message=r.content)
        return r.json()


class Followers(Resource):
    host = os.environ.get('SOCIAL_GRAPH_MANAGER_SERVICE_HOST', 'localhost')
    port = os.environ.get('SOCIAL_GRAPH_MANAGER_SERVICE_PORT', '9090')
    base_url = 'http://{}:{}'.format(host, port)

    def get(self):
        """Get users that follow current user
        """
        username, email = _get_user()
        r = requests.get('{}/followers/{}'.format(self.base_url, username))

        if not r.ok:
            abort(r.status_code, message=r.content)

        return r.json()

    def post(self):
        """Add a new follower to the current user"""
        username, email = _get_user()
        parser = RequestParser()
        parser.add_argument('follower', type=str, required=True)
        args = parser.parse_args()
        args.update(followed=username)
        r = requests.post(self.base_url + '/follow', json=args)
        if not r.ok:
            abort(r.status_code, message=r.content)

        return r.json()

class Following(Resource):
    host = os.environ.get('SOCIAL_GRAPH_MANAGER_SERVICE_HOST', 'localhost')
    port = os.environ.get('SOCIAL_GRAPH_MANAGER_SERVICE_PORT', '9090')
    base_url = 'http://{}:{}'.format(host, port)

    def get(self):
        """Get user current user is following
        """
        username, email = _get_user()
        r = requests.get('{}/following/{}'.format(self.base_url, username))

        if not r.ok:
            abort(r.status_code, message=r.content)

        return r.json()

    def post(self):
        """Have current user follow another user"""
        username, email = _get_user()
        parser = RequestParser()
        parser.add_argument('followed', type=str, required=True)
        args = parser.parse_args()
        args.update(follower=username)
        r = requests.post(self.base_url + '/follow', json=args)
        if not r.ok:
            abort(r.status_code, message=r.content)

        return r.json()

    def delete(self):
        """Have current user follow another user"""
        username, email = _get_user()
        parser = RequestParser()
        parser.add_argument('followed', type=str, required=True)
        args = parser.parse_args()
        args.update(follower=username)
        r = requests.post(self.base_url + '/unfollow', json=args)
        if not r.ok:
            abort(r.status_code, message=r.content)

        return r.json()
