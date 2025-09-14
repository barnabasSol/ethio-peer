using ResourceService.Models;
using Microsoft.EntityFrameworkCore;

namespace ResourceService.Repositories;

public class CourseRepo
{
    private readonly Context _context;

    public CourseRepo(Context context)
    {
        _context = context;
    }

    public async Task<Course?> GetCourseAsync(string courseCode)
    {
        return await _context.Courses.FindAsync(courseCode);
    }

    public async Task<IEnumerable<Course>> GetAllCoursesAsync()
    {
        return await _context.Courses.ToListAsync();
    }

    public async Task AddCourseAsync(CourseDTO course)
    {
        try
        {
            if (course == null)
            {
                throw new ArgumentNullException(nameof(course), "Course cannot be null");
            }
            var newCourse = new Course
            {
                CourseCode = course.CourseCode,
                Name = course.Name,
                Description = course.Description,
                CreditHour = course.CreditHour,
                Category = course.Category
            };
            await _context.Courses.AddAsync(newCourse);
            await _context.SaveChangesAsync();
        }
        catch
        {
            throw;
        }
    }
    public async Task<Course> UpdateCourseAsync(string courseCode, CourseDTO course)
    {
        try
        {
            if (courseCode.ToLower() != course.CourseCode.ToLower())
            {
                throw new ArgumentException("Course code mismatch");
            }

            var existingCourse = await _context.Courses.FindAsync(courseCode);
            if (existingCourse == null)
            {
                throw new KeyNotFoundException("Course not found");
            }
            existingCourse.Name = course.Name ?? existingCourse.Name;
            existingCourse.Description = course.Description ?? existingCourse.Description;
            existingCourse.CreditHour = course.CreditHour != 0 ?
            course.CreditHour : existingCourse.CreditHour;
            existingCourse.Category = course.Category.ToString() == "" ?
            existingCourse.Category : course.Category;

            _context.Courses.Update(existingCourse);
            await _context.SaveChangesAsync();
            return existingCourse;

        }
        catch
        {
            throw;
        }

    }
    public async Task DeleteCourseAsync(string courseCode)
    {
        try
        {
            var course = await _context.Courses.FindAsync(courseCode);
            if (course == null)
            {
                throw new KeyNotFoundException("Course not found");
            }
            _context.Courses.Remove(course);
            await _context.SaveChangesAsync();
        }
        catch
        {
            throw;
        }
    }

    public List<Course> GetCoursesByCategoryAsync(CourseCategory category)
    {
        var courses = _context.Courses.Where(c => c.Category == category).ToList();
        return courses;

    }
}