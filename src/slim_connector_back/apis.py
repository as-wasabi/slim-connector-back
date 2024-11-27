from slim_connector_back import app


@app.get("/health")
async def health():
    return {"ok": True}
