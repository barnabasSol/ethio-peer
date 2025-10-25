using Microsoft.EntityFrameworkCore;
using ResourceService.Models;
using ResourceService.Models.Dtos;

namespace ResourceService.Repositories;

public class RoomRepository(Context context)
{
    private readonly Context _context = context;
    //Add Room
    public async Task<Room> AddRoom(RoomDTO room)
    {
        var duplicateRoom = _context.Rooms.Where(r => r.SessionId == room.SessionId).FirstOrDefault();
        if (room == null) throw new ArgumentNullException(nameof(room), "Please provide a valid room");
        else if (duplicateRoom != null) throw new ArgumentException("Room with the same SessionId already exists");
        var newRoom = new Room
        {
            Id = Guid.NewGuid(),
            SessionId = room.SessionId,
            Name = room.Name,
            TopicId = room.TopicId
        };
        _context.Rooms.Add(newRoom);
        await _context.SaveChangesAsync();
        return newRoom;
    }
    //Delete Room
    public async Task DeleteRoom(Guid id)
    {
        var room = _context.Rooms.Where(r => r.Id == id).FirstOrDefault();
        if (room == null) throw new ArgumentNullException(nameof(id), "Room not found");
        _context.Rooms.Remove(room);
        await _context.SaveChangesAsync();
        return;
    }
    //Get Rooms
    public async Task<List<RoomResp>> GetRoomsAsync(CourseCategory? category = null, int count = 0)
    {
        var roomsQuery = _context.Rooms.Include(r => r.Topic)
           .Where(r => r.Topic != null).AsQueryable();
        if (category != null)
        {
            roomsQuery = roomsQuery
            .Where(r => r.Topic!.Course != null && r.Topic.Course.Category == category);
        }
        if (count > 0)
        {
            roomsQuery = roomsQuery.Take(count);
        }
        var roomsList = await roomsQuery.ToListAsync();

        return roomsList.Select(r => new RoomResp
        {
            RoomId = r.Id,
            RoomName = r.Name,
            TopicName = r.Topic!.Name,
            MemberCount = _context.RoomMembers.AsNoTracking().Count(m => m.RoomId == r.Id),
            CourseCode = r.Topic.CourseCode
        }).OrderByDescending(x => x.MemberCount)
            .ThenBy(x => x.RoomName)
            .ToList();
    }

    //Get Room by Id
    public async Task<Room> GetRoomById(Guid id)
    {
        var room = await _context.Rooms.Where(r => r.Id == id).FirstOrDefaultAsync();
        if (room == null) throw new ArgumentNullException(nameof(id), "Room not found");
        return room;
    }

    //Get Rooms by memberId
    public async Task<List<RoomResp>> GetRoomsByMemberId(string memberId, bool explore = false, int count = 0)
    {
        var roomsQuery = _context.Rooms
      .AsNoTracking()
      .Where(r => r.Topic != null);
        if (explore)
        {
            roomsQuery = roomsQuery.Where(r => !_context.RoomMembers.Any(m => m.RoomId == r.Id && m.UserId == memberId));
        }
        else
            roomsQuery = roomsQuery.Where(r => _context.RoomMembers.Any(m => m.RoomId == r.Id && m.UserId == memberId));
        if (count > 0)
        {
            roomsQuery = roomsQuery.Take(count);
        }
        var rooms = await roomsQuery.Include(r => r.Topic).ToListAsync();
        if (rooms.Count == 0) return [];
        var roomIds = rooms.Select(r => r.Id);
        var memberCounts = await _context.RoomMembers
    .AsNoTracking()
    .Where(m => roomIds.Contains(m.RoomId))
    .GroupBy(m => m.RoomId)
    .Select(g => new { RoomId = g.Key, Count = g.Count() })
    .ToListAsync();

        var roomMap = memberCounts.ToDictionary(x => x.RoomId, x => x.Count);


        return rooms.Select(r => new RoomResp
        {
            RoomId = r.Id,
            RoomName = r!.Name,
            TopicName = r.Topic!.Name,
            MemberCount = roomMap.TryGetValue(r.Id, out count) ? count : 0,
            CourseCode = r.Topic.CourseCode
        }).OrderByDescending(x => x.MemberCount)
                .ThenBy(x => x.RoomName)
                .ToList();
    }

    //Add member
    public async Task AddParticipant(RoomMemberDTO dto)
    {
        var room = _context.Rooms.Where(r => r.Id == dto.RoomId).FirstOrDefault() ?? throw new ArgumentNullException(nameof(dto.RoomId), "There is no Room with the given id");
        if (_context.RoomMembers.Where(r => r.RoomId == dto.RoomId && r.UserId == dto.MemberId).FirstOrDefault() != null) throw new ArgumentException("Member already in the room");
        _context.RoomMembers.Add(new RoomMember { RoomId = dto.RoomId, UserId = dto.MemberId });
        await _context.SaveChangesAsync();
    }
    //Remove member
    public async Task RemoveParticipant(RoomMemberDTO dto)
    {
        var room = _context.Rooms.Where(r => r.Id == dto.RoomId).FirstOrDefault() ?? throw new ArgumentNullException(nameof(dto.RoomId), "There is no Room with the given id");
        var member = _context.RoomMembers.Where(r => r.RoomId == dto.RoomId && r.UserId == dto.MemberId).FirstOrDefault() ?? throw new ArgumentNullException(nameof(dto.MemberId), "Member not found in the room");
        _context.RoomMembers.Remove(member);
        await _context.SaveChangesAsync();
    }
    //Get Members
    public async Task<List<string>> GetMembers(Guid roomId)
    {
        var room = _context.Rooms.Where(r => r.Id == roomId).FirstOrDefault() ?? throw new ArgumentNullException(nameof(roomId), "There is no Room with the given id");
        List<string> members = await _context.RoomMembers.AsNoTracking().Where(r => r.RoomId == roomId).Select(r => r.UserId).ToListAsync();
        return members;
    }

    public async Task<Guid> GetRoomIdBySessionId(string sessionId)
    {
        var roomId = await _context.Rooms.AsNoTracking().Where(r => r.SessionId == sessionId).Select(r => r.Id).FirstOrDefaultAsync();
        if (roomId == Guid.Empty) throw new ArgumentNullException(nameof(sessionId), "Room not found for the given sessionId");
        return roomId;

    }

    public async Task<List<RoomResp>> GetSuggestedRooms(List<string>? courses, string memberId)
    {
        if (courses == null || courses.Count == 0)
        {
            return await GetRoomsByMemberId(memberId, explore: true, count: 3);
        }

        var roomList = await _context.Rooms.AsNoTracking()
      .Where(r => r.Topic != null && r.Topic.Course != null && courses.Contains(r.Topic.Course.Name)
      && !_context.RoomMembers.Any(m => m.RoomId == r.Id && m.UserId == memberId))
      .Include(r => r.Topic)
      .Take(3)
      .ToListAsync();

        if (roomList == null || roomList.Count == 0)
            return await GetRoomsByMemberId(memberId, explore: true, count: 3);

        var roomIds = roomList.Select(r => r.Id);
        var memberCounts = await _context.RoomMembers
            .AsNoTracking()
            .Where(m => roomIds.Contains(m.RoomId))
            .GroupBy(m => m.RoomId)
            .Select(g => new { RoomId = g.Key, Count = g.Count() })
            .ToListAsync();

        var roomMap = memberCounts.ToDictionary(x => x.RoomId, x => x.Count);

        var rooms = roomList
            .Select(r => new RoomResp
            {
                RoomId = r.Id,
                TopicName = r.Topic!.Name,
                RoomName = r.Name,
                MemberCount = roomMap.TryGetValue(r.Id,out int count)?count:0,
                CourseCode = r.Topic.CourseCode
            })
            .OrderByDescending(x => x.MemberCount)
            .ThenBy(x => x.RoomName)
            .ToList();
        return rooms;
    }

}