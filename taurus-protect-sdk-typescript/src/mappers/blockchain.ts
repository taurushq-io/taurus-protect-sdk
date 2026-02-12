/**
 * Blockchain mapper functions for converting OpenAPI DTOs to domain models.
 */

import type {
  Blockchain,
  DOTBlockchainInfo,
  EVMBlockchainInfo,
  XTZBlockchainInfo,
} from '../models/blockchain';
import { currencyFromDto } from './currency';
import { safeBool, safeInt, safeMap, safeString } from './base';

/**
 * Maps a DOT blockchain info DTO to domain model.
 */
function dotBlockchainInfoFromDto(dto: unknown): DOTBlockchainInfo | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    ss58Format: safeInt(d.ss58Format ?? d.ss_58_format),
  };
}

/**
 * Maps an EVM blockchain info DTO to domain model.
 */
function evmBlockchainInfoFromDto(dto: unknown): EVMBlockchainInfo | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    chainId: safeString(d.chainId ?? d.chain_id),
  };
}

/**
 * Maps an XTZ blockchain info DTO to domain model.
 */
function xtzBlockchainInfoFromDto(dto: unknown): XTZBlockchainInfo | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    protocolHash: safeString(d.protocolHash ?? d.protocol_hash),
  };
}

/**
 * Maps a blockchain DTO to a Blockchain domain model.
 */
export function blockchainFromDto(dto: unknown): Blockchain | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    symbol: safeString(d.symbol),
    name: safeString(d.name),
    network: safeString(d.network),
    chainId: safeString(d.chainId ?? d.chain_id),
    confirmations: safeString(d.confirmations),
    blockHeight: safeString(d.blockHeight ?? d.block_height),
    blackholeAddress: safeString(d.blackholeAddress ?? d.blackhole_address),
    isLayer2Chain: safeBool(d.isLayer2Chain ?? d.is_layer2_chain),
    layer1Network: safeString(d.layer1Network ?? d.layer1_network),
    baseCurrency: currencyFromDto(d.baseCurrency ?? d.base_currency),
    dotInfo: dotBlockchainInfoFromDto(d.dotInfo ?? d.dot_info),
    ethInfo: evmBlockchainInfoFromDto(d.ethInfo ?? d.eth_info),
    xtzInfo: xtzBlockchainInfoFromDto(d.xtzInfo ?? d.xtz_info),
  };
}

/**
 * Maps an array of blockchain DTOs to Blockchain domain models.
 */
export function blockchainsFromDto(dtos: unknown[] | null | undefined): Blockchain[] {
  return safeMap(dtos, blockchainFromDto);
}
