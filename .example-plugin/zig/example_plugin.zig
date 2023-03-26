const std = @import("std");

pub fn range(len: usize) []const u0 {
    return @as([*]u0, undefined)[0..len];
}

export fn lpm_entrypoint(config_path: [*:0]u8, db_path: [*:0]u8, argc: usize, argv: [*][*:0]u8) void {
    std.debug.print("config_path: {s}\n", .{config_path});
    std.debug.print("db_path: {s}\n", .{db_path});

    for (range(argc)) |_, i| {
        std.debug.print("db_path: {s}\n", .{argv[i]});
    }
}
