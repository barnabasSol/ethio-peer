using Microsoft.EntityFrameworkCore;
using ResourceService.Models;

namespace ResourceService.Repositories;

public class RoomRepository(Context context)
{
    private readonly Context _context = context;
    //Add Room
    public async Task<Room> AddRoom(RoomDTO room)
    {
        var duplicateRoom = _context.Rooms.Where(r => r.SessionId.ToString() == room.SessionId.ToString()).FirstOrDefault();
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
    public async Task<List<Room>> GetRoomsByMemberId(Guid memberId)
    {
        List<Room?> rooms = await _context.RoomMembers.Where(r => r.UserId == memberId).Select(r => r.Room).ToListAsync();
        if (rooms == null || rooms.Count == 0) throw new ArgumentNullException(nameof(memberId), "No rooms found for the given userId");
        return rooms!;
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
    public async Task<List<Guid>> GetMembers(Guid roomId)
    {
        var room = _context.Rooms.Where(r => r.Id == roomId).FirstOrDefault() ?? throw new ArgumentNullException(nameof(roomId), "There is no Room with the given id");
        List<Guid> members = await _context.RoomMembers.Where(r => r.RoomId == roomId).Select(r => r.UserId).ToListAsync();
        return members;
    }

    public async Task<Guid> GetRoomIdBySessionId(Guid sessionId)
    {
        var roomId = await _context.Rooms.Where(r => r.SessionId == sessionId).Select(r => r.Id).FirstOrDefaultAsync();
        if (roomId == Guid.Empty) throw new ArgumentNullException(nameof(sessionId), "Room not found for the given sessionId");
        return roomId;

    }

    public async Task<List<RoomResp>> GetSuggestedRooms(List<string> courses)
    {
        if (courses == null || courses.Count == 0)
        {
            return await GetRoomsAsync(count: 3);
        }

        var roomList = await _context.Rooms.AsNoTracking()
      .Where(r => r.Topic != null && r.Topic.Course != null && courses.Contains(r.Topic.Course.Name))
      .Include(r => r.Topic)
      .Take(3)
      .ToListAsync();

        if (roomList == null || roomList.Count == 0)
            return await GetRoomsAsync(count: 3);

        var rooms = roomList
            .Select(r => new RoomResp
            {
                TopicName = r.Topic!.Name,
                RoomName = r.Name,
                MemberCount = _context.RoomMembers.AsNoTracking().Count(m => m.RoomId == r.Id),
                CourseCode = r.Topic.CourseCode
            })
            .OrderByDescending(x => x.MemberCount)
            .ThenBy(x => x.RoomName)
            .ToList();
        return rooms;
    }

}