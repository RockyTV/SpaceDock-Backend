from flask import request
from sqlalchemy import desc
from SpaceDock.routing import route
from SpaceDock.objects import *
from SpaceDock.common import *
from SpaceDock.formatting import *

import json


@route('/api/access')
@user_has('access-view')
def roles():
    """
    Displays  list of all roles with the matching abilities
    """
    roles = list()
    for role in Role.query.order_by(desc(Role.id)).all():
        roles.append({'id': role.id, 'name': role.name, 'abilities': bulk(role.abilities, ability_format), 'params': json.loads(role.params)})
    return {'error': False, 'count': len(roles), 'data': roles}

@route('/api/access/roles/assign', methods=['POST'])
@user_has('access-edit')
@with_session
def add_role():
    """
    Promotes a user for the given role. Required parameters: userid, rolename
    """
    userid = request.form.get('userid')
    rolename = request.form.get('rolename')
    if not userid.isdigit() or not User.query.filter(User.id == int(userid)).first():
        return {'error': True, 'reasons': ['The userid is invalid.']}, 400
    user = User.query.filter(User.id == int(userid)).first()
    user.add_roles(rolename)
    return {'error': False}

@route('/api/access/roles/remove', methods=['POST'])
@user_has('access-edit')
@with_session
def remove_role():
    """
    Removes a user from a group. Required parameters: userid, rolename
    """
    userid = request.form.get('userid')
    rolename = request.form.get('rolename')
    if not userid.isdigit() or not User.query.filter(User.id == int(userid)).first():
        return {'error': True, 'reasons': ['The userid is invalid.']}, 400
    user = User.query.filter(User.id == int(userid)).first()
    if not rolename in user._roles:
        return {'error': True, 'reasons': ['The user doesn\'t have this role']}, 400
    user.remove_roles(rolename)
    return {'error': False}

@route('/api/access/abilities/assign', methods=['POST'])
@user_has('access-edit')
@with_session
def add_abilities():
    """
    Adds a permission to a group. Required parameters: rolename, abname
    """
    rolename = request.form.get('rolename')
    abname = request.form.get('abname')
    if not Role.query.filter(Role.name == rolename).first():
        return {'error': True, 'reasons': ['The role does not exist. Please add it to a user to create it internally.']}, 400
    r = Role.query.filter(Role.name == rolename).first()
    r.add_abilities(abname)
    return {'error': False}

@route('/api/access/abilities/remove', methods=['POST'])
@user_has('access-edit')
@with_session
def remove_abilities():
    """
    Removes a permission from a group. Required parameters: rolename, abname
    """
    rolename = request.form.get('rolename')
    abname = request.form.get('abname')
    errors = []
    if not Role.query.filter(Role.name == rolename).first():
        errors.append('The role does not exist.')
    if not Ability.query.filter(Ability.name == abname).first():
        errors.append('The ability does not exist.')
    if len(errors) > 0:
        return {'error': True, 'reasons': errors}, 400
    role = Role.query.filter(Role.name == rolename).first()
    ability = Ability.query.filter(Ability.name == abname).first()
    if not ability in role.abilities:
        return {'error': True, 'reasons': ['The ability isn\'t assigned to this role']}, 400
    role.remove_abilities(abname)
    return {'error': False}

@route('/api/access/params/add/<rolename>', methods=['POST'])
@user_has('access-edit')
@with_session
def add_params(rolename):
    """
    Adds a parameter for an ability. Required parameters: abname, param, value
    """
    abname = request.form.get('abname')
    param = request.form.get('param')
    value = request.form.get('value')
    errors = []
    if not Role.query.filter(Role.name == rolename).first():
        errors.append('The rolename is invalid.')
    if not Ability.query.filter(Ability.name == abname).first():
        errors.append('The ability does not exist.')
    if len(errors) > 0:
        return {'error': True, 'reasons': errors}, 400
    role = Role.query.filter(Role.name == rolename).first()
    role.add_param(abname, param, value)
    return {'error': False}

@route('/api/access/params/remove/<rolename>', methods=['POST'])
@user_has('access-edit')
@with_session
def remove_params(rolename):
    """
    Removes a parameter from an ability. Required parameters: abname, param, value
    """
    abname = request.form.get('abname')
    param = request.form.get('param')
    value = request.form.get('value')
    errors = []
    if not Role.query.filter(Role.name == rolename).first():
        errors.append('The rolename is invalid.')
    if not Ability.query.filter(Ability.name == abname).first():
        errors.append('The ability does not exist.')
    if len(errors) > 0:
        return {'error': True, 'reasons': errors}, 400
    role = Role.query.filter(Role.name == rolename).first()
    if not value in get_param(abname, param, json.loads(role.params)):
        return {'error': True, 'reasons': ['The parameter doesn\'t exist']}, 400
    role.remove_param(abname, param, value)
    return {'error': False}