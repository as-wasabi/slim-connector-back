import pytest
import pytest_asyncio
import sqlalchemy

from slim_connector_back import tbls, mdls, ENV
from slim_connector_back.user.__res import SelfUserRes
from slim_connector_back.user.__user_body import PostUserBody
from slim_connector_back.util import keycloak, tks
from test.conftest import session


@pytest.fixture
def post_user_body(session) -> PostUserBody:
    return PostUserBody(
        user_name="PostUserBody_user_name",
        user_icon_uuid=None,
    )


@pytest_asyncio.fixture
async def newuser_keycloak_profile(session) -> keycloak.KeycloakUserProfile:
    uid = "917ebffb-0e86-4189-87d1-604f7246be29"
    return keycloak.KeycloakUserProfile(
        sub=uid,
        email_verified=True,
        preferred_username="newuser_keycloak_profile",
        email="newuser_keycloak_profile@example.com",
    )


@pytest.fixture
def newuser_access_token(session, newuser_keycloak_profile) -> tks.TokenInfo:
    return mdls.JwtTokenData.new(
        mdls.TokenType.access,
        newuser_keycloak_profile
    ).new_token_info(ENV.token.secret_key)


@pytest.mark.asyncio
async def test_create_user(session, client, newuser_access_token, post_user_body):
    result = await client.post(
        "/api/user",
        post_user_body,
        newuser_access_token.token
    )
    assert result.status_code == 200, f"invalid status code {result.json()}"
    body = result.json()
    assert body is not None
    body = SelfUserRes(**body)
    result = await session.execute(
        sqlalchemy.select(sqlalchemy.func.count())
        .select_from(tbls.UserTable)
        .where(tbls.UserTable.user_id == body.user_id)
    )
    assert result.scalar_one() == 1, f"\n{body}\n"


@pytest.mark.asyncio
async def test_get_self(client, login_access_token, session, login_user, login_keycloak_profile):
    result = await client.get(
        "/api/user/self",
        login_access_token.token
    )
    assert result.status_code == 200, f"invalid status code {result.read()}"
    body = result.json()
    assert body is not None
    body = SelfUserRes(**body)
    assert body.user_id == login_keycloak_profile.sub
    assert body.user_mail == login_keycloak_profile.email
    assert body.user_name == login_user.user_name
