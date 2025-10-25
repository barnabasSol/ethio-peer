using ResourceService.Models;
using ResourceService.Models.Dtos;
using ResourceService.Repositories;
using Microsoft.AspNetCore.Mvc;
namespace ResourceService.Controllers;
[ApiController]
[Route("[controller]s")]
public class TopicController(TopicRepo topicRepo) : ControllerBase
{
    public required TopicRepo _topicRepo = topicRepo;

    [HttpGet]
    public async Task<IActionResult> GetTopics()
    {
        var topics = await _topicRepo.GetAllTopicsAsync();
        return Ok(topics);
    }
    
    [HttpGet("{id}")]
    public async Task<IActionResult> GetTopic(Guid id)
    {
        Topic topic = await _topicRepo.GetTopicAsync(id);
        return Ok(new
        {
            topic.TopicId,
            topic.Name,
            topic.Description,
            topic.CourseCode
        });
    }
    [HttpGet("pattern")]
    public async Task<IActionResult> GetTopicsByPattern([FromQuery] string pattern = "")
    {
        return Ok(await _topicRepo.GetTopicsByPattern(pattern));
    }
    [HttpGet("chart")]
    public async Task<IActionResult> GetTopicVsDocument()
    {
        return Ok(new { items=await _topicRepo.GetTopicDocCount() });
    }

    
    [HttpPost]
    public async Task<IActionResult> CreateTopic([FromBody] TopicDTO topic)
    {
        try
        {
            Topic newTopic=await _topicRepo.AddTopicAsync(topic);
            return CreatedAtAction(nameof(CreateTopic), newTopic);
        }
        catch (ArgumentNullException ex)
        {
            return BadRequest(ex.Message);
        }
        catch
        {
            throw;
        }
    }
    [HttpPut("{id}")]
    public async Task<IActionResult> UpdateTopic(Guid id, [FromBody] TopicDTO topic)
    {
        try
        {
            Topic updatedTopic = await _topicRepo.UpdateTopicAsync(id, topic);
            return Ok(new
            {
                updatedTopic.TopicId,
                updatedTopic.CourseCode,
                updatedTopic.Name,
                updatedTopic.Description,
                updatedTopic.Course
            });
        }
        catch (KeyNotFoundException ex)
        {
            return NotFound(ex.Message);
        }
        catch (ArgumentException ex)
        {
            return BadRequest(ex.Message);
        }
        catch (Exception ex)
        {
            return StatusCode(500, $"Internal server error: {ex.Message}");
        }
    }

    [HttpDelete("{id}")]
    public async Task<IActionResult> DeleteTopic(Guid id)
    {
        try
        {
            await _topicRepo.DeleteTopicAsync(id);
            return NoContent();
        }
        catch (KeyNotFoundException ex)
        {
            return NotFound(ex.Message);
        }
        catch (Exception ex)
        {
            return StatusCode(500, $"Internal server error: {ex.Message}");
        }
    }

} 