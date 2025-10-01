using ResourceService.Models;
using ResourceService.Repositories;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
namespace ResourceService.Controllers;
[ApiController]
[Route("api/[controller]s")]
public class CourseController(CourseRepo courseRepo) : ControllerBase
{
    public required CourseRepo _courseRepo = courseRepo;

    [HttpGet]
    public async Task<IActionResult> GetCourses()
    {
        var course = await _courseRepo.GetAllCoursesAsync();
        return Ok(course);
    }

    [HttpGet("{code}")]
    public async Task<IActionResult> GetCourse(string code)
    {
        Course ?course = await _courseRepo.GetCourseAsync(code);
        if (course == null)
        {
            return NotFound($"Course with code {code} not found");
        }
        return Ok(new
        {
            course.CourseCode,
            course.Name,
            course.Description,
            course.CreditHour
        });
    }

    [HttpGet("category/{category}")]
    public IActionResult GetCoursesByCategory(CourseCategory category)
    {
        var courses = _courseRepo.GetCoursesByCategoryAsync(category);
        return Ok(courses);
    }

    [HttpPost]
    public async Task<IActionResult> CreateCourse([FromBody] CourseDTO course)
    {
        try
        {
            await _courseRepo.AddCourseAsync(course);
            return CreatedAtAction(nameof(CreateCourse), course);
        }
        catch (ArgumentNullException ex)
        {
            return BadRequest(ex.Message);
        }
        catch (Exception ex)
        {
            return StatusCode(500, $"Internal server error: {ex.Message}");
        }
    }
    [HttpPut("{code}")]
    public async Task<IActionResult> UpdateCourse(string code, [FromBody] CourseDTO course)
    {
        try
        {
            Course updatedCourse = await _courseRepo.UpdateCourseAsync(code, course);
            return Ok(new
            {
                updatedCourse.CourseCode,
                updatedCourse.Name,
                updatedCourse.Description,
                updatedCourse.CreditHour
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

    [HttpDelete("{code}")]
    public async Task<IActionResult> DeleteCourse(string code)
    {
        try
        {
            await _courseRepo.DeleteCourseAsync(code);
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