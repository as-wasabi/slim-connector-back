import pytest
import pytest_asyncio
import sqlalchemy

from slim_connector_back import tbls
from slim_connector_back.chat.__body import PostChatBody, PostChatMessageBody
from slim_connector_back.chat.__res import ChatMessageRes, ChatRes, ChatMessagesRes
from slim_connector_back.chat.__result import ChatMessageResult


@pytest.fixture
def post_chat_body(session, saved_user) -> PostChatBody:
    return PostChatBody(
        users=[
            saved_user.user_id
        ]
    )


@pytest.fixture
def post_chat_message_body(session) -> PostChatMessageBody:
    return PostChatMessageBody(
        message="post_chat_message_body",
        images=[]
    )


@pytest_asyncio.fixture
async def saved_chat(session, saved_user, login_user_deps) -> ChatRes:
    body = PostChatBody(
        users=[
            saved_user.user_id
        ]
    )
    res = await body.save_new(login_user_deps, session)
    return res.to_chat_res()


@pytest_asyncio.fixture
async def saved_message(session, login_user_deps, saved_chat) -> ChatMessageResult:
    body = PostChatMessageBody(
        message="saved_message",
        images=[],
    )
    await login_user_deps.refresh(session)
    res = await body.save_new(session, saved_chat.chat_id, login_user_deps)
    await session.commit()
    return res


@pytest.fixture
def saved_chat_list(saved_chat) -> list[ChatRes]:
    return [saved_chat]


@pytest.fixture
def saved_message_list(saved_message) -> list[ChatMessageResult]:
    return [saved_message]


@pytest.mark.asyncio
async def test_post_chat(session, client, login_access_token, post_chat_body):
    response = await client.post(
        "/api/chat",
        post_chat_body,
        login_access_token.token
    )
    assert response.status_code == 200, f"invalid status code {response.json()}"
    body = response.json()
    assert body is not None
    chat = ChatRes(**body)

    record = await session.execute(
        sqlalchemy.select(tbls.ChatUserTable)
        .where(tbls.ChatUserTable.chat_id == chat.chat_id)
    )
    user_tbls: list[tbls.ChatUserTable] = record.scalars().all()
    assert len(user_tbls) == len(chat.users)
    for i in range(len(user_tbls)):
        user_tbl = user_tbls[i]
        user_id = chat.users[i]
        assert user_tbl.chat_id == chat.chat_id
        assert user_tbl.user_id == user_id


@pytest.mark.asyncio
async def test_get_chat(client, login_access_token, session, saved_chat_list):
    result = await client.get(
        "/api/chat",
        login_access_token.token
    )
    assert result.status_code == 200, f"invalid status code {result.read()}"
    body = result.json()
    assert body is not None
    assert len(saved_chat_list) == len(body)
    for i in range(len(saved_chat_list)):
        chat = saved_chat_list[i]
        res = ChatRes(**body[i])
        assert chat.chat_id == res.chat_id
        for user_i in range(len(chat.users)):
            assert chat.users[user_i] == res.users[user_i]


@pytest.mark.asyncio
async def test_get_message(client, login_access_token, session, saved_chat, saved_message_list):
    result = await client.get(
        f"/api/chat/{saved_chat.chat_id}/message",
        login_access_token.token
    )
    assert result.status_code == 200, f"invalid status code {result.read()}"
    body = result.json()
    assert body is not None
    res = ChatMessagesRes(**body)
    assert len(saved_message_list) == len(res.messages)
    assert saved_chat.chat_id == res.chat_id
    for i in range(len(saved_message_list)):
        message_record = saved_message_list[i]
        message_res = res.messages[i]
        await message_record.refresh(session)
        assert message_record.message.chat_message_id == message_res.chat_message_id
        assert message_record.message.message == message_res.message
        assert message_record.message.post_user_id == message_res.post_user_id
        for j in range(len(message_record.images)):
            assert message_record.images[j].image_uuid == message_res.images[j]


@pytest.mark.asyncio
async def test_post_message(
        session,
        client,
        login_access_token,
        saved_chat,
        post_chat_message_body,
        login_user,
):
    response = await client.post(
        f"/api/chat/{saved_chat.chat_id}/message",
        post_chat_message_body,
        login_access_token.token
    )
    assert response.status_code == 200, f"invalid status code {response.json()}"
    body = response.json()
    assert body is not None
    message = ChatMessageRes(**body)

    result = await session.execute(
        sqlalchemy.select(tbls.ChatMessageTable)
        .where(tbls.ChatMessageTable.chat_message_id == message.chat_message_id)
    )
    record: tbls.ChatMessageTable = result.scalar_one()

    assert record.chat_id == message.chat_id
    assert record.chat_message_id == message.chat_message_id
    assert record.message == message.message
    assert record.index == message.index
    assert record.post_user_id == login_user.user_id

    result = await session.execute(
        sqlalchemy.select(tbls.ChatImageTable)
        .where(tbls.ChatImageTable.chat_message_id == record.chat_id)
    )
    image_records: list[tbls.ChatImageTable] = result.scalars().all()

    assert len(image_records) == len(message.images)
    for i in range(len(image_records)):
        image = image_records[i]
        assert image.image_uuid == message.images[i]
