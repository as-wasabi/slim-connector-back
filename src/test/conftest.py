import asyncio
from typing import Iterator

import pytest
import pytest_asyncio
from _pytest.fixtures import FixtureRequest
from sqlalchemy import NullPool
from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker

from slim_connector_back import main, ENV, deps, mdls
from slim_connector_back.db import BaseTable
from slim_connector_back.deps import JwtTokenDeps, UserDeps
from slim_connector_back.user.__res import SelfUserRes
from slim_connector_back.user.__user_body import PostUserBody
from slim_connector_back.util import keycloak, tks
from test.base import Client


@pytest.fixture(scope="session")
def event_loop(request: FixtureRequest) -> Iterator[asyncio.AbstractEventLoop]:
    loop = asyncio.get_event_loop_policy().new_event_loop()
    loop.__original_fixture_loop = True  # type: ignore[attr-defined]
    yield loop
    loop.close()


@pytest.fixture(scope="session")
def app():
    return main.app


@pytest.fixture(scope="session")
def engine(event_loop):
    engine = create_async_engine(ENV.database.db_url, echo=False, poolclass=NullPool)
    try:
        yield engine
    finally:
        engine.sync_engine.dispose()


@pytest_asyncio.fixture
async def create(engine):
    try:
        async with engine.begin() as conn:
            await conn.run_sync(BaseTable.metadata.create_all)
        yield
    finally:
        async with engine.begin() as conn:
            await conn.run_sync(BaseTable.metadata.drop_all)


@pytest_asyncio.fixture
async def session(engine, create, app):
    session_maker = async_sessionmaker(autocommit=False, autoflush=False, bind=engine)

    async def override_get_db():
        async with session_maker() as db_session:
            yield db_session
            await db_session.commit()

    app.dependency_overrides[deps.DbDeps.session] = override_get_db

    async with session_maker() as session:
        yield session


@pytest.fixture
def client(app):
    return Client(app)


@pytest_asyncio.fixture
async def login_keycloak_profile(session) -> keycloak.KeycloakUserProfile:
    uid = "565dc2fd-127d-45d1-9a10-8002280777d7"
    return keycloak.KeycloakUserProfile(
        sub=uid,
        email_verified=True,
        preferred_username="user_login",
        email="user_login@example.com",
    )


@pytest_asyncio.fixture
async def login_user(session, login_keycloak_profile) -> SelfUserRes:
    body = PostUserBody(
        user_name="user_login",
        user_icon_uuid=None,
    )
    res = await body.save_new(session, login_keycloak_profile)
    return res.to_self_user_res()


@pytest.fixture
def login_access_token(session, login_keycloak_profile, login_user) -> tks.TokenInfo:
    return mdls.JwtTokenData.new(
        mdls.TokenType.access,
        login_keycloak_profile
    ).new_token_info(ENV.token.secret_key)


@pytest_asyncio.fixture
async def login_access_jwt_token_deps(session, login_user, login_access_token) -> JwtTokenDeps | None:
    return deps.JwtTokenDeps.get_token_or_none(login_access_token.token)


@pytest_asyncio.fixture
async def login_user_deps(session, login_user, login_access_jwt_token_deps) -> UserDeps:
    return await deps.UserDeps.get(session, login_access_jwt_token_deps)


@pytest_asyncio.fixture
async def saved_user(session) -> SelfUserRes:
    uid = "df56c011-8025-468a-a390-202e6f0d6328"
    profile = keycloak.KeycloakUserProfile(
        sub=uid,
        email_verified=True,
        preferred_username="post_chat_body",
        email="post_chat_body@example.com",
    )
    body = PostUserBody(
        user_name="PostUserBody_user_name_saved",
        user_icon_uuid=None,
    )
    res = await body.save_new(session, profile)
    return res.to_self_user_res()
