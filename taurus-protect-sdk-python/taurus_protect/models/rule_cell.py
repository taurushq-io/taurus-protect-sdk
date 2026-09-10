"""Typed transaction-rule cell union (cross-SDK contract).

Each cell type corresponds to a column type and one of its cell-type enum
values in the governance protobuf schema (e.g. ``FiatAmountRange`` is the
``RuleFiatAmountRange`` cell of a ``RuleFiatAmount`` column). The ``kind``
discriminator carries the model type name.

The ``*Any`` cells are cells like any other; being the protobuf zero values,
their serialized form is the empty cell. Cells whose type/content this SDK
version cannot represent losslessly are preserved verbatim as :class:`RawCell`,
so nothing is ever silently dropped. The same set exists in the Go, Java and
TypeScript SDKs and the wire encoding is pinned by the shared golden vectors at
``scripts/resources/governance-cell-vectors.json``.
"""

from __future__ import annotations

from typing import List, Literal, Union

from pydantic import BaseModel, Field
from typing_extensions import Annotated

_CFG = {"frozen": True}


# --- RuleFiatAmount column ---
class FiatAmountAny(BaseModel):
    kind: Literal["FiatAmountAny"] = "FiatAmountAny"
    model_config = _CFG


class FiatAmountIsZero(BaseModel):
    kind: Literal["FiatAmountIsZero"] = "FiatAmountIsZero"
    model_config = _CFG


class FiatAmountRange(BaseModel):
    kind: Literal["FiatAmountRange"] = "FiatAmountRange"
    min_amount: str = ""
    max_amount: str = ""
    model_config = _CFG


# --- RuleSource column ---
class SourceAny(BaseModel):
    kind: Literal["SourceAny"] = "SourceAny"
    model_config = _CFG


class SourceInternalWallet(BaseModel):
    kind: Literal["SourceInternalWallet"] = "SourceInternalWallet"
    path: str = ""
    model_config = _CFG


class SourceInternalAddress(BaseModel):
    kind: Literal["SourceInternalAddress"] = "SourceInternalAddress"
    address: str = ""
    path: str = ""
    model_config = _CFG


class SourceAnyExchange(BaseModel):
    kind: Literal["SourceAnyExchange"] = "SourceAnyExchange"
    model_config = _CFG


class SourceExchange(BaseModel):
    kind: Literal["SourceExchange"] = "SourceExchange"
    label: str = ""
    model_config = _CFG


class SourceExternalAddress(BaseModel):
    kind: Literal["SourceExternalAddress"] = "SourceExternalAddress"
    address: str = ""
    memo: str = ""
    model_config = _CFG


# --- RuleDestination column ---
class DestinationAny(BaseModel):
    kind: Literal["DestinationAny"] = "DestinationAny"
    model_config = _CFG


class DestinationInternalWallet(BaseModel):
    kind: Literal["DestinationInternalWallet"] = "DestinationInternalWallet"
    path: str = ""
    model_config = _CFG


class DestinationInternalAddress(BaseModel):
    kind: Literal["DestinationInternalAddress"] = "DestinationInternalAddress"
    address: str = ""
    path: str = ""
    model_config = _CFG


class DestinationExternalAddress(BaseModel):
    kind: Literal["DestinationExternalAddress"] = "DestinationExternalAddress"
    address: str = ""
    memo: str = ""
    model_config = _CFG


class DestinationAnyExchange(BaseModel):
    kind: Literal["DestinationAnyExchange"] = "DestinationAnyExchange"
    model_config = _CFG


class DestinationExchange(BaseModel):
    kind: Literal["DestinationExchange"] = "DestinationExchange"
    label: str = ""
    memo: str = ""
    model_config = _CFG


class DestinationContractAddress(BaseModel):
    kind: Literal["DestinationContractAddress"] = "DestinationContractAddress"
    address: str = ""
    name: str = ""
    symbol: str = ""
    blockchain: str = ""
    model_config = _CFG


class DestinationAnyExternalAddress(BaseModel):
    kind: Literal["DestinationAnyExternalAddress"] = "DestinationAnyExternalAddress"
    model_config = _CFG


class DestinationAnyContractAddress(BaseModel):
    kind: Literal["DestinationAnyContractAddress"] = "DestinationAnyContractAddress"
    model_config = _CFG


# --- RuleStringEqual column (payload is raw UTF-8 string bytes) ---
class StringEqualAny(BaseModel):
    kind: Literal["StringEqualAny"] = "StringEqualAny"
    model_config = _CFG


class StringEqualEmpty(BaseModel):
    kind: Literal["StringEqualEmpty"] = "StringEqualEmpty"
    model_config = _CFG


class StringEqualValue(BaseModel):
    kind: Literal["StringEqualValue"] = "StringEqualValue"
    value: str = ""
    model_config = _CFG


# --- RuleBytesEqual column (payload is raw bytes) ---
class BytesEqualAny(BaseModel):
    kind: Literal["BytesEqualAny"] = "BytesEqualAny"
    model_config = _CFG


class BytesEqualEmpty(BaseModel):
    kind: Literal["BytesEqualEmpty"] = "BytesEqualEmpty"
    model_config = _CFG


class BytesEqualValue(BaseModel):
    kind: Literal["BytesEqualValue"] = "BytesEqualValue"
    value: bytes = b""
    model_config = _CFG


# --- RuleStringArrayEqual column ---
class StringArrayEqualAny(BaseModel):
    kind: Literal["StringArrayEqualAny"] = "StringArrayEqualAny"
    model_config = _CFG


class StringArrayEqualEmpty(BaseModel):
    kind: Literal["StringArrayEqualEmpty"] = "StringArrayEqualEmpty"
    model_config = _CFG


class StringArrayEqualValue(BaseModel):
    kind: Literal["StringArrayEqualValue"] = "StringArrayEqualValue"
    values: List[str] = Field(default_factory=list)
    model_config = _CFG


# --- RuleIntegerGreater column (sign selects the Value/NegValue wire arm) ---
class IntegerGreaterAny(BaseModel):
    kind: Literal["IntegerGreaterAny"] = "IntegerGreaterAny"
    model_config = _CFG


class IntegerGreaterValue(BaseModel):
    kind: Literal["IntegerGreaterValue"] = "IntegerGreaterValue"
    value: int = 0
    model_config = _CFG


# --- RuleUIntegerGreater column ---
class UIntegerGreaterAny(BaseModel):
    kind: Literal["UIntegerGreaterAny"] = "UIntegerGreaterAny"
    model_config = _CFG


class UIntegerGreaterIsZero(BaseModel):
    kind: Literal["UIntegerGreaterIsZero"] = "UIntegerGreaterIsZero"
    model_config = _CFG


class UIntegerGreaterValue(BaseModel):
    kind: Literal["UIntegerGreaterValue"] = "UIntegerGreaterValue"
    value: int = 0
    model_config = _CFG


class UIntegerGreaterIsEqual(BaseModel):
    kind: Literal["UIntegerGreaterIsEqual"] = "UIntegerGreaterIsEqual"
    value: int = 0
    model_config = _CFG


# --- RuleWhitelistedContract column ---
class WhitelistedContractAny(BaseModel):
    kind: Literal["WhitelistedContractAny"] = "WhitelistedContractAny"
    model_config = _CFG


class WhitelistedContractAddress(BaseModel):
    kind: Literal["WhitelistedContractAddress"] = "WhitelistedContractAddress"
    address: str = ""
    name: str = ""
    symbol: str = ""
    blockchain: str = ""
    model_config = _CFG


# --- Fallback for cells unknown to this SDK version (preserved verbatim) ---
class RawCell(BaseModel):
    kind: Literal["RawCell"] = "RawCell"
    column_type: str = ""
    payload: bytes = b""
    model_config = _CFG


RuleCell = Annotated[
    Union[
        FiatAmountAny, FiatAmountIsZero, FiatAmountRange,
        SourceAny, SourceInternalWallet, SourceInternalAddress, SourceAnyExchange, SourceExchange, SourceExternalAddress,
        DestinationAny, DestinationInternalWallet, DestinationInternalAddress, DestinationExternalAddress,
        DestinationAnyExchange, DestinationExchange, DestinationContractAddress, DestinationAnyExternalAddress, DestinationAnyContractAddress,
        StringEqualAny, StringEqualEmpty, StringEqualValue,
        BytesEqualAny, BytesEqualEmpty, BytesEqualValue,
        StringArrayEqualAny, StringArrayEqualEmpty, StringArrayEqualValue,
        IntegerGreaterAny, IntegerGreaterValue,
        UIntegerGreaterAny, UIntegerGreaterIsZero, UIntegerGreaterValue, UIntegerGreaterIsEqual,
        WhitelistedContractAny, WhitelistedContractAddress,
        RawCell,
    ],
    Field(discriminator="kind"),
]


_FAMILY = {
    "FiatAmountAny": "RuleFiatAmount", "FiatAmountIsZero": "RuleFiatAmount", "FiatAmountRange": "RuleFiatAmount",
    "SourceAny": "RuleSource", "SourceInternalWallet": "RuleSource", "SourceInternalAddress": "RuleSource",
    "SourceAnyExchange": "RuleSource", "SourceExchange": "RuleSource", "SourceExternalAddress": "RuleSource",
    "DestinationAny": "RuleDestination", "DestinationInternalWallet": "RuleDestination",
    "DestinationInternalAddress": "RuleDestination", "DestinationExternalAddress": "RuleDestination",
    "DestinationAnyExchange": "RuleDestination", "DestinationExchange": "RuleDestination",
    "DestinationContractAddress": "RuleDestination", "DestinationAnyExternalAddress": "RuleDestination",
    "DestinationAnyContractAddress": "RuleDestination",
    "StringEqualAny": "RuleStringEqual", "StringEqualEmpty": "RuleStringEqual", "StringEqualValue": "RuleStringEqual",
    "BytesEqualAny": "RuleBytesEqual", "BytesEqualEmpty": "RuleBytesEqual", "BytesEqualValue": "RuleBytesEqual",
    "StringArrayEqualAny": "RuleStringArrayEqual", "StringArrayEqualEmpty": "RuleStringArrayEqual",
    "StringArrayEqualValue": "RuleStringArrayEqual",
    "IntegerGreaterAny": "RuleIntegerGreater", "IntegerGreaterValue": "RuleIntegerGreater",
    "UIntegerGreaterAny": "RuleUIntegerGreater", "UIntegerGreaterIsZero": "RuleUIntegerGreater",
    "UIntegerGreaterValue": "RuleUIntegerGreater", "UIntegerGreaterIsEqual": "RuleUIntegerGreater",
    "WhitelistedContractAny": "RuleWhitelistedContract", "WhitelistedContractAddress": "RuleWhitelistedContract",
}


def cell_family(cell: BaseModel) -> str:
    """The column type a typed cell belongs to (empty for RawCell -> its column_type)."""
    if cell.kind == "RawCell":
        return cell.column_type  # type: ignore[attr-defined]
    return _FAMILY.get(cell.kind, "")
