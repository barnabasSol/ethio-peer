using Microsoft.AspNetCore.Mvc;
using ResourceService.Models;
using ResourceService.Repositories;

namespace ResourceService.Controllers;

[ApiController]
[Route("[controller]s")]
public class RoomController(RoomRepository roomRepo) : ControllerBase
{
    private readonly RoomRepository _roomRepo = roomRepo;
    //Get Rooms
    [HttpGet]
    public async Task<IActionResult> GetRooms([FromQuery] Guid? memberId,CourseCategory? category)
    {
        try
        {
            if (memberId.HasValue && memberId.Value != Guid.Empty)
            {
                var rooms = await _roomRepo.GetRoomsByMemberId(memberId.Value);
                return Ok(rooms);
            }
            else
            {
                var rooms = await _roomRepo.GetRoomsAsync(category);
                return Ok(rooms);
            }
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
    [HttpGet("suggestions")]
    public async Task<IActionResult> GetSuggestedRooms([FromQuery] List<string> courses)
    { 
        var rooms = await _roomRepo.GetSuggestedRooms(courses);
        return Ok(rooms);
    }


    //Get Rooms by Id
    [HttpGet("{id}")]
    public async Task<IActionResult> GetRoomById(Guid id)
    {
        try
        {
            var room = await _roomRepo.GetRoomById(id);
            return Ok(room);
        }
        catch (ArgumentNullException e)
        {
            return NotFound(e.ParamName);
        }
        catch (Exception ex)
        {
            return StatusCode(500, ex.Message);
        }
    }

    // //Get Rooms by memberId
    // [HttpGet]
    // public async Task<IActionResult> GetRoomsByMemberId([FromQuery]Guid memberId)
    // {
    //     try
    //     {
    //         var rooms = await _roomRepo.GetRoomsByMemberId(memberId);
    //         return Ok(rooms);
    //     }
    //     catch (ArgumentNullException e)
    //     {
    //         return NotFound(e.Message);
    //     }
    //     catch (Exception ex)
    //     {
    //         return StatusCode(500, ex.Message);
    //     }
    // }

    // Add Room
    [HttpPost]
    public async Task<IActionResult> AddRoom([FromBody] RoomDTO room)
    {
        try
        {
            var newRoom = await _roomRepo.AddRoom(room);
            return CreatedAtAction(nameof(AddRoom), newRoom);
        }
        catch (ArgumentNullException e)
        {
            return BadRequest(e.ParamName);
        }
        catch (ArgumentException e)
        {
            return Conflict(e.Message);
        }
        catch (Exception ex)
        {
            return StatusCode(500, ex.Message);
        }
    }
    [HttpDelete("{id}")]
    // Delete Room
    public async Task<IActionResult> DeleteRoom(Guid id)
    {
        try
        {
            await _roomRepo.DeleteRoom(id);
            return NoContent();
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
    //// Add Participant
    [HttpPost("members")]
    public async Task<IActionResult> AddParticipant([FromBody] RoomMemberDTO roomMember)
    {
       try
       {
           if (roomMember.MemberId == Guid.Empty || roomMember.RoomId == Guid.Empty) return BadRequest("Please provide a valid parameter");
           await _roomRepo.AddParticipant(roomMember);
           return NoContent();
       }
       catch (ArgumentNullException e)
       {
           return NotFound(e.Message);
       }
       catch (ArgumentException e)
       {
           return Conflict(e.Message);
       }
       catch (Exception ex)
       {
           return StatusCode(500, ex.Message);
       }
    }
    // Remove Participant
    [HttpDelete("members")]
    public async Task<IActionResult> RemoveParticipant(RoomMemberDTO dto)
    {
       try
       {
           await _roomRepo.RemoveParticipant(dto);
           return NoContent();
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
    //returns a list of members in a room
    [HttpGet("{roomId}/members")]
    public async Task<IActionResult> GetMembersInRoom(Guid roomId)
    {
       try
       {
           var members = await _roomRepo.GetMembers(roomId);
           return Ok(members);
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