import logging

import uvicorn
from fastapi import Depends
from fastapi import FastAPI

from slim_connector_back.service import EnvService, DbEnvService, CorsService
from slim_connector_back.util import fastapiutil

logger = logging.getLogger(__name__)
logger.addHandler(logging.StreamHandler())
logger.setLevel("INFO")

app = FastAPI(dependencies=[
    Depends(EnvService),
    Depends(DbEnvService),
    Depends(CorsService),
])


# logging.getLogger("uvicorn.access").addFilter(ExcludeFilter(["/health"]))


class Main:
    def __init__(self, fast_api: FastAPI):
        self.app = fast_api
        fastapiutil.handler(fast_api)
        # noinspection PyUnresolvedReferences
        import slim_connector_back.apis

    def main(self):
        uvicorn.run(self.app, host="0.0.0.0")


main = Main(app)
