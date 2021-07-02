USE [minesweeper]
GO
/****** Object:  StoredProcedure [dbo].[CreateGame]    Script Date: 7/1/2021 05:29:40 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE PROCEDURE [dbo].[CreateSpot] (
    @gameId INT,
    @value VARCHAR(20),
	@x INT,
	@y INT,
	@nearSpots VARCHAR(50),
	@status VARCHAR(20),
	@id INT OUTPUT
) AS
BEGIN
    INSERT INTO [dbo].[Spot] ([GameId],[Value],[X],[Y],[NearSpots],[Status])
     VALUES (@gameId, @value, @x, @y, @nearSpots, @status)

    SET @id=SCOPE_IDENTITY()
    RETURN  @id
END;