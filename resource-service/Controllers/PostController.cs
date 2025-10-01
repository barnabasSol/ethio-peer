using Microsoft.AspNetCore.Mvc;
using ResourceService.Repositories;

namespace ResourceService.Controllers;

[ApiController]
[Route("api/[controller]s")]
public class PostController(PostRepo postRepo) : ControllerBase
{
    private readonly PostRepo _postRepo = postRepo;

    // Get Posts by RoomId
    [HttpGet]
    public async Task<IActionResult> GetPostsByRoomId([FromQuery] Guid roomId)
    {
        try
        {
            var posts = await _postRepo.GetPostsByRoomId(roomId);
            return Ok(posts);
        }
        catch (ArgumentNullException e)
        {
            return NotFound(e.Message);
        }
        catch (Exception ex)
        {
            return StatusCode(500, ex.Message);
        }
    }

    
}  