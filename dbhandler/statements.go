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
const SELECT_GAMES_BY_USER = `
	SELECT [GameId],[UserId],[CreatedDate],[TimeConsumed],[Status],[Rows],[Columns],[Mines] 
		FROM [minesweeper].[dbo].[Game] 
	WHERE UserId = @p1 AND STATUS='Pending'`

const SELECT_GAME_BY_ID = `
	SELECT [GameId],[UserId],[CreatedDate],[TimeConsumed],[Status],[Rows],[Columns],[Mines] 
		FROM [minesweeper].[dbo].[Game] 
	WHERE GameId = @p1`

const SELECT_LATEST_GAME = `
	SELECT TOP 1 [GameId],[UserId],[CreatedDate],[TimeConsumed],[Status],[Rows],[Columns],[Mines] 
		FROM [minesweeper].[dbo].[Game] 
	WHERE UserId = @p1 ORDER BY CreatedDate DESC`

const SELECT_SPOTS_BY_GAME_ID = `
	SELECT TOP (1000) [SpotId],[GameId],[Value],[X],[Y],[NearSpots],[Status]
		FROM [minesweeper].[dbo].[Spot] 
	WHERE GameId = @p1`
