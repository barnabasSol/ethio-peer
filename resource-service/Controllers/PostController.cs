using Microsoft.AspNetCore.Mvc;
using ResourceService.Models;
using ResourceService.Repositories;

namespace ResourceService.Controllers;

[ApiController]
[Route("[controller]s")]
public class PostController(PostRepo postRepo) : ControllerBase
{
    private readonly PostRepo _postRepo = postRepo;

    // Get Posts by RoomId
    [HttpGet]
    public async Task<IActionResult> GetPostsByRoomId([FromQuery] Guid roomId, [FromQuery] PagedQuery pq)
    {
        try
        {
            var posts = await _postRepo.GetPostsByRoomId(roomId, pq);
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
    [HttpGet("weekly")]
    public async Task<IActionResult> GetWeeklyPosts()
    {
        return Ok(await _postRepo.GetWeeklyPosts());
    }

}
