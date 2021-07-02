USE [minesweeper]
GO
/****** Object:  StoredProcedure [dbo].[CreateGame]    Script Date: 7/2/2021 02:21:44 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
ALTER PROCEDURE [dbo].[CreateGame] (
    @userId INT,
    @timeConsumed INT,
	@rows INT,
	@columns INT,
	@mines INT,
	@id INT OUTPUT
) AS
BEGIN
    INSERT INTO [dbo].[Game] ([UserId],[CreatedDate],[TimeConsumed],[Status],[Rows],[Columns],[Mines]) 
	VALUES 
	(@userId, GetDate(), @timeConsumed, 'Pending', @rows, @columns, @mines)

    SET @id=SCOPE_IDENTITY()
    RETURN  @id
END;