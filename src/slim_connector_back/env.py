import os


from slim_connector_back.util import urls



class Token:
    refresh_token_expire_minutes: int | float
    access_token_expire_minutes: int | float
    algorithm: str

    def __init__(self):
        refresh_token_expire_minutes = os.getenv("REFRESH_TOKEN_EXPIRE_MINUTES")
        if refresh_token_expire_minutes is None:
            self.refresh_token_expire_minutes = 60 * 24 * 14
        else:
            self.refresh_token_expire_minutes = float(refresh_token_expire_minutes)

        access_token_expire_minutes = os.getenv("ACCESS_TOKEN_EXPIRE_MINUTES")
        if access_token_expire_minutes is None:
            self.access_token_expire_minutes = 15
        else:
            self.access_token_expire_minutes = float(access_token_expire_minutes)

        algorithm = os.getenv("ALGORITHM")
        if algorithm is None:
            self.algorithm = "HS256"
        else:
            self.algorithm = algorithm

        secret_key = os.getenv("SECRET_KEY")
        if secret_key is None:
            raise ValueError
        else:
            self.secret_key = secret_key


