# Differ fixtures

`known-gap/` is a four-SDK surface set with one service deliberately missing from one
SDK. `selftest.sh` runs the differ against it and fails if the differ reports parity —
without this, a broken extractor (an empty or truncated surface) would make the differ
print "OK" and the gate would be worthless while looking green.
