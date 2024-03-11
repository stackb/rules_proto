from module_app.app import api_pb2
from util import date_pb2

request = api_pb2.ApiRequest(
    resource_names = ["foo", "bar", "baz"],
    date = date_pb2.Date(
        year = "2023",
        month = "12",
        day = "20",
    ),
)

print(request)
