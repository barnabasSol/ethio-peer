using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.AspNetCore.Mvc;
using ResourceService.Models;
using ResourceService.Repositories;

namespace ResourceService.Controllers;

[ApiController]
[Route("api/[controller]s")]
public class DocumentController(DocRepo docRepo) : ControllerBase
{
    private readonly DocRepo _docRepo = docRepo;
    [HttpGet]
    public async Task<IActionResult> GetDocs()
    {
        var docs = await _docRepo.GetDocsAsync();
        return Ok(docs);
    }

    //upload doc
    [HttpPost]
    [IgnoreAntiforgeryToken]
    public async Task<IActionResult> DocumentUpload([FromForm] DocDTO dto)
    {
        try
        {
            var doc = await _docRepo.AddDoc(dto);
            return CreatedAtAction(nameof(DocumentUpload), doc);
        }
        catch (FileNotFoundException)
        {
            return NotFound();
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
    //download doc
    [HttpGet("{id}")]
    public async Task<IActionResult> DocumentDownload(string id)
    {
        try
        {
            var stream = await _docRepo.DownloadDoc(id);
            return File(stream, "application/octet-stream", id);
        }
        catch (FileNotFoundException)
        {
            return NotFound();
        }
    }

    //update metadata
    [HttpPut("{id}")]
    public async Task<IActionResult> DocumentUpdate(string id, [FromForm] DocDTO dto)
    {
        try
        {
            var updatedDoc = await _docRepo.ModifyDocMetadata(id, dto);
            return Ok(updatedDoc);
        }
        catch (FileNotFoundException)
        {
            return NotFound();
        }
        catch (ArgumentNullException e)
        {
            return BadRequest(e.ParamName);
        }
        catch (Exception ex)
        {
            return StatusCode(500, ex.Message);
        }
    }
    //delete doc
    [HttpDelete("{id}")]
    public async Task<IActionResult> DeleteDoc(string id)
    {
        try
        {
            var deleted = await _docRepo.DeleteDoc(id);
            return deleted ? NoContent() : NotFound();
        }
        catch (FileNotFoundException)
        {
            return NotFound();
        }
    }
    //get doc metadata by room id
    // [HttpGet]
    //  public async Task<IActionResult> GetDocsByRoomId([FromQuery] Guid roomId)
    // {
    //     try

    //     {
    //         DateTime.Now.ToString();
    //             var docs = await _docRepo.GetDocsByRoomId(roomId);
    //         return Ok(docs);
    //     }
    //     catch (Exception ex)
    //     {
    //         return StatusCode(500, ex.Message);
    //     }
    // }
    //get doc metadata by topic id
}