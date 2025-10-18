using Microsoft.EntityFrameworkCore;

namespace ResourceService.Models;

public class PagedList<T>
{
    public PagedList(int pageSize, int pageNumber, int totalCount, List<T> items)
    {
        Items = items;
        PageSize = pageSize;
        PageNumber = pageNumber;
        TotalCount = totalCount;
    }
    public List<T> Items { get; set; } = [];
    public int TotalCount { get; set; } = 0;
    public int PageSize { get; set; } = 0;
    public int PageNumber { get; set; } = 0;
    public bool HasNextPage => TotalCount > PageNumber * PageSize;
    public int ItemCount => Items.Count;
    public bool HasPreviousPage => PageNumber > 1;
    public int MaxPageCount => TotalCount == 0 || PageSize == 0 ? 0 : (int)Math.Ceiling(TotalCount / (double)PageSize);
    public async static Task<PagedList<T>> CreateAsync(IQueryable<T> query, int? pageSize, int? pageNumber)
    {
        
        int count = await query.CountAsync();
        var currentPage = pageNumber ?? 1;
        var itemsPerPage = pageSize ?? count; 
        var items = await query.Skip((currentPage - 1) * itemsPerPage).Take(itemsPerPage).ToListAsync();

        return new PagedList<T>(
            itemsPerPage,
           currentPage,
          count, items
        );

    }

}

