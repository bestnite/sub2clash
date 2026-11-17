import assert from "node:assert/strict";
import { test } from "node:test";
import { appUrl } from "../src/utils.ts";

for (const prefix of ["/", "/sub2clash/", "/tools/sub2clash/"]) {
  for (const trailingSlash of [true, false]) {
    const pathname = trailingSlash ? prefix : prefix.replace(/\/$/, "");
    const page = `https://example.com${pathname}?theme=dark#section`;
    test(`application URLs preserve ${pathname || "/"}`, () => {
      for (const path of [
        "short",
        "short/example?password=secret",
        "convert/eyJ0ZXN0Ijp0cnVlfQ==",
        "s/example?password=secret",
      ]) {
        assert.equal(appUrl(path, page), `https://example.com${prefix}${path}`);
      }
    });
  }
}

test("short-link passwords retain reserved characters", () => {
  const password = "a&b#c+d? e";
  const link = appUrl(
    `s/example?${new URLSearchParams({ password })}`,
    "https://example.com/sub2clash/"
  );
  assert.equal(new URL(link).searchParams.get("password"), password);
  assert.equal(new URL(link).pathname, "/sub2clash/s/example");
});
