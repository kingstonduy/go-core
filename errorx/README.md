#### USAGE

##### 1. Default errors

| Error Name                | Status (int)                  | Error Code (string) | Message                      |
| ------------------------- | ----------------------------- | ------------------- | ---------------------------- |
| `Success`               | `200 OK`                    | `00`               | Successful                   |
| `Failed`                | `200 OK`                    | `01`       | Failed. %s                                 |
| `ValidationError`       | `200 OK`                    | `02`       | Validation error. %s                       |
| `NotFoundError`         | `200 OK`                    | `03`       | Not found error. %s                        |
| `OutboundError`         | `200 OK`                    | `04`       | Outbound error. %s                         |
| `TimeoutError`          | `200 OK`                    | `05`       | Timeout error. %s                          |
| `BadRequestError`       | `200 OK`                    | `06`       | Bad request error. %s                      |
| `UnauthorizedError`     | `200 OK`                    | `07`       | Unauthorized error. %s                     |
| `ForbiddenError`        | `200 OK`                    | `08`       | Forbidden error. %s                        |
| `MethodNotAllowedError` | `200 OK`                    | `09`       | Method not allowed error. %s               |
| `ConflictError`         | `200 OK`                    | `10`       | Conflict error. %s                         |
| `TooManyRequestError`   | `200 OK`                    | `11`       | Too many request error. %s                 |
| `InternalServerError`   | `500 Internal Server Error` | `999`      | Internal server error. %s                  |

##### 2. Create dynamic func error

````go
	var (
		UserNotActive = NewError(400, "54", "User not active")
	)

	func UserNotFound(userId int) error {
		return NewError(400, "54", "User not found %d", userId)
	}
````
