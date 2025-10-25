using ResourceService.Models;
using Microsoft.EntityFrameworkCore;
using ResourceService.Models.Dtos;

namespace ResourceService.Repositories;

public class TopicRepo
{
    private readonly Context _context;

    public TopicRepo(Context context)
    {
        _context = context;
    }

    public async Task<Topic> GetTopicAsync(Guid id)
    {
        try
        {
            Topic? topic = await _context.Topics.FindAsync(id);
            if (topic == null)
            {
                throw new KeyNotFoundException("Topic not found");
            }
            return topic;
        }
        catch
        {
            throw;
        }
    }

    public async Task<IEnumerable<Topic>> GetAllTopicsAsync()
    {
        return await _context.Topics.ToListAsync();
    }
    public async Task<List<TopicMin>> GetTopicsByPattern(string pattern)
    {
        return await _context.Topics.Where(t => EF.Functions.ILike(t.Name, $"{pattern}%")).Select(t => new TopicMin
        {
            TopicId = t.TopicId,
            TopicName = t.Name
        }).ToListAsync();
    }
    public async Task<List<TopicDocRes>> GetTopicDocCount()
    {
        var query = from d in _context.Documents
                    join r in _context.Rooms on d.RoomId equals r.Id
                    join t in _context.Topics on r.TopicId equals t.TopicId
                    group d by new { t.TopicId, t.Name } into g
                    select new TopicDocRes
                    {
                        TopicName = g.Key.Name,
                        DocCount = g.Count()
                    };
        return await query.ToListAsync();
    }

    public async Task<Topic> AddTopicAsync(TopicDTO topic)
    {
        try
        {
            Course? foundCourse = await _context.Courses.FindAsync(topic.CourseCode);
            if (topic == null)
            {
                throw new ArgumentNullException("Topic cannot be null");
            }
            else if (foundCourse == null)
            {
                throw new KeyNotFoundException("Course not found");
            }
            var newTopic = new Topic
            {
                TopicId = Guid.NewGuid(),
                CourseCode = topic.CourseCode,
                Name = topic.Name,
                Description = topic.Description
            };
            await _context.Topics.AddAsync(newTopic);
            await _context.SaveChangesAsync();
            return newTopic;
        }
        catch
        {
            throw;
        }
    }
    public async Task<Topic> UpdateTopicAsync(Guid id, TopicDTO topicDTO)
    {
        try
        {

            var existingTopic = await _context.Topics.FindAsync(id);
            if (existingTopic == null)
            {
                throw new KeyNotFoundException("Topic not found");
            }
            existingTopic.Name = topicDTO.Name ?? existingTopic.Name;
            existingTopic.Description = topicDTO.Description ?? existingTopic.Description;
            existingTopic.CourseCode = topicDTO.CourseCode ?? existingTopic.CourseCode;
            _context.Topics.Update(existingTopic);
            await _context.SaveChangesAsync();
            return existingTopic;

        }
        catch
        {
            throw;
        }

    }
    public async Task DeleteTopicAsync(Guid id)
    {
        try
        {
            var topic = await _context.Topics.FindAsync(id);
            if (topic == null)
            {
                throw new KeyNotFoundException("Topic not found");
            }
            _context.Topics.Remove(topic);
            await _context.SaveChangesAsync();
        }
        catch
        {
            throw;
        }
    }

    public async Task<string> GetTopicNameById(Guid id)
    {
        String topicName = await _context.Topics.Where(t => t.TopicId == id).Select(t => t.Name).FirstAsync();
        if (topicName == null)
        {
            throw new KeyNotFoundException("Topic not found");
        }
        return topicName;

    }
}