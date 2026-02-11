/**
 * Score models for Taurus-PROTECT SDK.
 *
 * Scores represent risk assessments provided by third-party blockchain analytics
 * providers (e.g., Chainalysis, Elliptic). These scores help assess the risk
 * associated with blockchain addresses for anti-money laundering (AML) and
 * compliance purposes.
 */

/**
 * Represents a compliance score from a risk assessment provider.
 *
 * Different providers use different scoring methodologies and scales.
 * The score value should be interpreted according to the specific
 * provider's documentation.
 *
 * @example
 * ```typescript
 * // Access scores from an address
 * const address = await client.addresses.get(addressId);
 * for (const score of address.scores ?? []) {
 *   console.log(`Provider: ${score.provider}`);
 *   console.log(`Score: ${score.score}`);
 *   console.log(`Type: ${score.type}`);
 * }
 * ```
 */
export interface Score {
  /** The unique identifier of the score record */
  readonly id?: number;

  /**
   * The name of the risk assessment provider.
   * Common providers include "chainalysis", "coinfirm", "elliptic", "scorechain".
   */
  readonly provider?: string;

  /**
   * The type of score.
   * Common types include "in" (for incoming transaction risk) and "out" (for outgoing).
   */
  readonly type?: string;

  /**
   * The score value as assigned by the provider.
   * The interpretation of this value depends on the provider's scoring methodology.
   */
  readonly score?: string;

  /** The date and time when the score was last updated */
  readonly updateDate?: Date;
}
