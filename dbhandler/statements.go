package dbhandler

const VALIDATE_LOGIN = "SELECT [UserId],[UserName],[UserLastName],[Password],[CreatedDate] FROM [dbo].[User] WHERE [UserName] = @p1 AND [Password] = @p2"

const CREATE_USER = `
DECLARE	@return_value int,
		@id int

EXEC	@return_value = [dbo].[CreateUser]
		@name = @p1,
		@lastName = @p2,
		@password = @p3,
		@id = @id OUTPUT

SELECT	@id as N'@id'
`
const CREATE_GAME = `
DECLARE	@return_value int,
		@id int

EXEC	@return_value = [dbo].[CreateGame]
		@userId = @p1,
		@timeConsumed = @p2,
		@status =@p3,
		@rows = @p4,
		@columns = @p5,
		@mines = @p6,
		@id = @id OUTPUT

SELECT	@id as N'@id'
`

const CREATE_SPOT = `
DECLARE	@return_value int,
		@id int

EXEC	@return_value = [dbo].[CreateSpot]
		@gameId = @p1,
		@value = @p2,
		@x = @p3,
		@y = @p4,
		@nearSpots = @p5,
		@status = @p6,
		@id = @id OUTPUT

SELECT	@id as N'@id'
`
