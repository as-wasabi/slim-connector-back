import uuid
from datetime import datetime

import pytest, pytest_asyncio

from slim_connector_back.tbls import CartTable, CartProductTable, ProductTable

from test.conftest import session, login_keycloak_profile, login_user_deps

from slim_connector_back.util.tks import TokenInfo

from slim_connector_back.responses import CartProduct




#ユーザーデータ
# @pytest_asyncio.fixture
# async def mock_user(session):
#     user_id = uuid.UUID("ad89b295-54fc-4b94-b8d4-e9e5850e27ce")
#     user = UserTable(
#         user_id=user_id,
#         user_name="Ado",
#         user_screen_id="d89c7adb-74b2-f904-322e-70f642ee8132",
#         user_icon_uuid=uuid.UUID("0325de47-2abe-d6d9-e99e-da1a0c3c1f3e"),
#         user_date=datetime(2024,11,9,22,15, 34),
#         user_mail="ado@example.com",
#     )
#     session.add(user)
#     await session.commit()
#     await session.refresh(user)
#     return user

# 商品データ
@pytest_asyncio.fixture
async def mock_product(session)-> ProductTable:
    product_id = uuid.UUID("835dc2fd-127d-45d1-9a10-8002280777d8")
    product = ProductTable(
        product_id=product_id,
        product_price=300,
        product_title="缶バッジ",
        product_text="ランダムで全20種類",
        product_date=datetime(2024,11,9,23,15, 34),
        product_contents_uuid=uuid.UUID("3c4b6ab2-b82f-3cad-0cf4-a3a6612b7236"),
        product_thumbnail_uuid=uuid.UUID("0325de47-2abe-d6d9-e99e-da1a0c3c1f3e"),
    )
    session.add(product)
    await session.commit()
    await session.refresh(product)
    return product

#  -> GetProductsResponse

# refreshは → データ取り出し
# flushは非同期じゃないと、できない
# add →　commit
# add →　flush
# commitするならflushは不要

# カートデータ
@pytest_asyncio.fixture
async def mock_cart(session, login_keycloak_profile)-> CartTable:
    cart_id = uuid.UUID("3a5bfd4a-a2cb-8914-d0df-a139c5176f85")
    user_id = login_keycloak_profile.sub
    print(f"user_id→→→{user_id}")

    cart = CartTable(
        cart_id=cart_id,
        user_id=user_id,
        purchase_date=None
    )
    session.add(cart)
    await session.commit()
    await session.refresh(cart)
    return cart


# カートプロダクトデータ
@pytest_asyncio.fixture
async def create_mock_cart_product(session, mock_cart, mock_product) -> CartProductTable:
    await session.refresh(mock_cart)
    await session.refresh(mock_product)
    cart_product = CartProductTable(
        cart_id=mock_cart.cart_id,
        product_id= mock_product.product_id,
    )
    session.add(cart_product)
    await session.commit()
    await session.refresh(cart_product)
    return cart_product


@pytest_asyncio.fixture
async def mock_data(session, mock_product, mock_cart, create_mock_cart_product):
    # await session.refresh(mock_user)
    await session.refresh(mock_product)
    await session.refresh(mock_cart)
    await session.refresh(create_mock_cart_product)

    return mock_product, mock_cart, create_mock_cart_product

@pytest.fixture
def cart_product_mock():
    """CartProduct モックデータ"""
    return CartProduct(
        product_id="a94cf31c-1e14-4189-9869-d1a455fb6529",
        product_price=1000,
        product_title="Test Product",
        product_text="This is a test product.",
        product_date=datetime(2022,12,9,23,15, 34),
        product_contents_uuid="31fa694d-f5fa-4e75-880c-9bce43bfbe39",
        product_thumbnail_uuid="ee3cb623-15a0-487b-960a-7a23f6e9fea9",
    )

# @pytest_asyncio.fixture
# async def put_product_cart(session) -> CartProduct:


@pytest.mark.asyncio
async def test_read_product_cart(
    client,
    session,
    login_access_token,
    # login_user_deps,
    mock_data
):


    # デバッグログを出力
    # product, cart, cart_product = mock_data

    # await session.refresh(user)
    # await session.refresh(product)
    # await session.refresh(cart)
    # await session.refresh(cart_product)
    #
    # # print(f"mock_user_id: {user.user_id}")
    # print(f"mock_cart_id: {cart.cart_id}")
    # print(f"mock_cart_product.cart_id: {cart_product.cart_id} mock_cart_product.product_id--->{cart_product.product_id}")
    # print(f"login_access_token.token---> {login_access_token.token}")


    # 非同期的にGETリクエストを送信
    response = await client.get(
        "/api/cart_product",
        login_access_token.token
    )
    assert response.status_code == 200, f"Unexpected status code: {response.status_code}, {response.json()}"


    # レスポンスデータ 検証
    data = response.json()
    print(f"data----->{data}")

    # データがリストであることを確認
    assert isinstance(data, list), "レスポンスデータはリストではありません"

    # リストの中身が辞書であることを確認
    assert isinstance(data[0], dict), "レスポンスデータの要素は辞書ではありません"

    # 必須キー 確認
    required_keys = {
        'product_id', 'product_title', 'product_text',
        'product_thumbnail_uuid', 'product_date',
        'product_price', 'product_contents_uuid'
    }
    assert required_keys <= data[0].keys(), "レスポンスデータに必須キーが不足しています"


    # 値の型 確認
    assert isinstance(data[0]['product_id'], str), "product_idは文字列ではありません"
    assert isinstance(data[0]['product_title'], str), "product_titleは文字列ではありません"
    assert isinstance(data[0]['product_text'], str), "product_textは文字列ではありません"
    assert isinstance(data[0]['product_thumbnail_uuid'], str), "product_thumbnail_uuidは文字列ではありません"
    assert isinstance(data[0]['product_date'], str), "product_dateは文字列ではありません"
    assert isinstance(data[0]['product_price'], int), "product_priceは整数ではありません"
    assert isinstance(data[0]['product_contents_uuid'], str), "product_contents_uuidは文字列ではありません"

    # 値の正確性を確認（価格の値範囲チェック）
    assert data[0]['product_price'] > 0, "product_priceが0以下です"



@pytest.mark.asyncio
async def test_put_cart_product(
        client,
        session,
        login_access_token,
        login_user_deps,
        cart_product_mock
):

    token_info = TokenInfo(
        token=login_access_token.token,
        expire=login_access_token.expire
    )
    print(f"token_info---> {token_info.token}")

    token_info_json = token_info.model_dump_json()
    print(f"token_info_json---> {token_info_json}")



    response = await client.put(
        "/api/cart_product",
        cart_product_mock,
        login_access_token.token
    )
    assert response.status_code == 200
    response_data = response.json()
    assert response_data == {}

