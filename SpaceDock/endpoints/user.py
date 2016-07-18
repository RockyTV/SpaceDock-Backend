from flask import request
from flask_login import current_user
from sqlalchemy import desc
from SpaceDock.common import *
from SpaceDock.formatting import user_info, admin_user_info
from SpaceDock.objects import *
from SpaceDock.routing import route

import json


@route('/api/users')
def get_users():
    """
    Returns a list of users
    """
    users = list()
    for user in User.query.order_by(desc(User.id)).all():
        if has_ability('view-users-full') or (user.public and user.confirmation == None):
            users.append(user_info(user) if not has_ability('view-users-full') else admin_user_info(user))
    return {'error': False, 'count': len(users), 'data': users}

@route('/api/users/<userid>')
def get_user_info(userid):
    """
    Returns more data for one user
    """
    user = None
    if userid == 'current':
        user = current_user
    elif not userid.isdigit() or not User.query.filter(User.id == int(userid)).first():
        return {'error': True, 'reasons': ['The userid is invalid']}, 400
    if not user:
        user = User.query.filter(User.id == int(userid)).first()
    if has_ability('view-users-full') or (user.public and user.confirmation == None):
        return {'error': False, 'count': 1, 'data': (user_info(user) if not has_ability('view-users-full') else admin_user_info(user))}
    else:
        return {'error': True, 'reasons': ['The userid is invalid']}, 400

@route('/api/users/<userid>/edit', methods=['POST'])
@user_has('user-edit', params=['userid'])
@with_session
def edit_user(userid):
    """
    Edits a user, based on the request parameters. Required fields: data
    """
    if not userid.isdigit() or not User.query.filter(User.id == int(userid)).first():
        return {'error': True, 'reasons': ['The userid is invalid.']}, 400

    # Get variables
    parameters = json.loads(request.form['data'])

    # Get the matching user and edit it
    user = User.query.filter(User.id == int(userid)).first()
    edit_object(user, parameters)
    return {'error': False}