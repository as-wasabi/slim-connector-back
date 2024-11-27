import uuid
from datetime import datetime

import pytest
import pytest_asyncio
import sqlalchemy

from slim_connector_back import tbls
from slim_connector_back.product.__res import GetProductsResponse
from test.conftest import session


@pytest_asyncio.fixture
async def product_table_saved(session) -> tbls.ProductTable:
    table = tbls.ProductTable.insert(
        session,
        product_price=100,
        product_title="title",
        product_text="text",
        product_date=datetime.now(),
        product_contents_uuid=uuid.uuid4(),
        product_thumbnail_uuid=uuid.uuid4(),
    )
    await session.commit()
    await session.refresh(table)
    return table


@pytest.mark.asyncio
async def test_read_products(client, session, product_table_saved):
    result = await client.get(
        "/products"
    )
    # ステータスコードの検証
    assert result.status_code == 200, f"invalid status code {result.read()}"

    # レスポンスボディの確認
    body = result.json()
    assert body is not None

    # 取得したデータの数と、APIのレスポンスに含まれるデータの数（bodyの長さ）が一致していることを確認
    records = await session.execute(
        sqlalchemy.select(tbls.ProductTable).where()
    )
    records = records.scalars().all()
    assert len(records) == len(body)

    # データの個々の検証
    for i in range(len(records)):
        record = records[i]
        product = GetProductsResponse(**body[i])

        assert record.product_id == product.product_id
        assert record.product_date == product.product_date
        assert record.product_text == product.product_text
        assert record.product_title == product.product_title
        assert record.product_price == product.product_price
