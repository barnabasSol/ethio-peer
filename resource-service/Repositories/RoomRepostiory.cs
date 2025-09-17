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
        if (room == null) throw new ArgumentNullException(nameof(room),"Please provide a valid room");
        else if (duplicateRoom != null) throw new ArgumentException("Room with the same SessionId already exists");
        var newRoom = new Room
        {
            Id = Guid.NewGuid(),
            SessionId = room.SessionId,
            Name = room.Name,
        };
        _context.Rooms.Add(newRoom);
        await _context.SaveChangesAsync();
        return newRoom;
    }
    //Delete Room
    public async Task DeleteRoom(Guid id)
    {
        var room = _context.Rooms.Where(r => r.Id == id).FirstOrDefault();
        if (room == null) throw new ArgumentNullException(nameof(id),"Room not found");
        _context.Rooms.Remove(room);
        await _context.SaveChangesAsync();
        return;
    }
    //Get Rooms
    public IEnumerable<Room> GetRooms()
    {
        return _context.Rooms;
    }
    //Get Room by Id
    public async Task<Room> GetRoomById(Guid id)
    {
        var room = await _context.Rooms.Where(r => r.Id == id).FirstOrDefaultAsync();
        if (room == null) throw new ArgumentNullException(nameof(id),"Room not found");
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
        var room = _context.Rooms.Where(r => r.Id == dto.RoomId).FirstOrDefault() ?? throw new ArgumentNullException(nameof(dto.RoomId),"There is no Room with the given id");
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


}