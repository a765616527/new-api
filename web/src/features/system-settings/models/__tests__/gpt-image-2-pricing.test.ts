/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { describe, expect, test } from 'vitest'

import {
  buildPreviewRows,
  EMPTY_LANE_ENABLED,
  EMPTY_LANE_PRICES,
  needsGPTImage2FallbackPrice,
} from '../model-pricing-core'
import { buildModelSnapshots } from '../model-pricing-snapshots'

const emptySnapshotInput = {
  modelRatio: '{}',
  cacheRatio: '{}',
  createCacheRatio: '{}',
  completionRatio: '{}',
  imageRatio: '{}',
  audioRatio: '{}',
  audioCompletionRatio: '{}',
  billingMode: '{}',
  billingExpr: '{}',
}

describe('GPT Image 2 resolution pricing', () => {
  test('attaches GPT Image 2.5 flare tier keys to the flare model', () => {
    const snapshots = buildModelSnapshots({
      ...emptySnapshotInput,
      modelPrice: JSON.stringify({
        'gpt-image-2.5-flare@1k': 0.03,
        'gpt-image-2.5-flare@4k': 0.08,
      }),
    })

    expect(snapshots).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          name: 'gpt-image-2.5-flare',
          billingMode: 'per-request',
          gptImage2Price1K: '0.03',
          gptImage2Price4K: '0.08',
        }),
      ])
    )
  })

  test('keeps tier keys internal and attaches them to the public model', () => {
    const snapshots = buildModelSnapshots({
      ...emptySnapshotInput,
      modelPrice: JSON.stringify({
        'gpt-image-2': 0.03,
        'gpt-image-2@1k': 0.01,
        'gpt-image-2@2k': 0.02,
        'gpt-image-2@4k': 0.04,
      }),
    })

    expect(snapshots).toHaveLength(1)
    expect(snapshots[0]).toMatchObject({
      name: 'gpt-image-2',
      billingMode: 'per-request',
      price: '0.03',
      gptImage2Price1K: '0.01',
      gptImage2Price2K: '0.02',
      gptImage2Price4K: '0.04',
    })
  })

  test('restores an incomplete tier-only configuration so it can be repaired', () => {
    const snapshots = buildModelSnapshots({
      ...emptySnapshotInput,
      modelPrice: JSON.stringify({ 'gpt-image-2@1k': 0.01 }),
    })

    expect(snapshots).toHaveLength(1)
    expect(snapshots[0]).toMatchObject({
      name: 'gpt-image-2',
      billingMode: 'per-request',
      price: '',
      gptImage2Price1K: '0.01',
    })
  })

  test('requires a fallback unless every resolution tier has a price', () => {
    expect(
      needsGPTImage2FallbackPrice({
        name: 'gpt-image-2',
        price: '',
        gptImage2Price1K: '0.01',
        gptImage2Price2K: '',
        gptImage2Price4K: '0.04',
      })
    ).toBe(true)
    expect(
      needsGPTImage2FallbackPrice({
        name: 'gpt-image-2',
        price: '',
        gptImage2Price1K: '0.01',
        gptImage2Price2K: '0.02',
        gptImage2Price4K: '0.04',
      })
    ).toBe(false)
    expect(
      needsGPTImage2FallbackPrice({
        name: 'gpt-image-2',
        price: '0.03',
        gptImage2Price1K: '0.01',
        gptImage2Price2K: '',
        gptImage2Price4K: '',
      })
    ).toBe(false)
    expect(
      needsGPTImage2FallbackPrice({
        name: 'gpt-image-2.5-flare',
        price: '',
        gptImage2Price1K: '0.03',
        gptImage2Price2K: '',
        gptImage2Price4K: '0.08',
      })
    ).toBe(true)
  })

  test('previews configured tiers and marks missing tiers as fallback', () => {
    const rows = buildPreviewRows(
      {
        name: 'gpt-image-2',
        price: '0.03',
        gptImage2Price1K: '0.01',
        gptImage2Price2K: '',
        gptImage2Price4K: '0.04',
      },
      'per-request',
      '',
      '',
      '',
      EMPTY_LANE_PRICES,
      EMPTY_LANE_ENABLED,
      (key) => key
    )

    expect(rows).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ key: 'gpt-image-2-1k', value: '$0.01' }),
        expect.objectContaining({ key: 'gpt-image-2-2k', value: 'Fallback' }),
        expect.objectContaining({ key: 'gpt-image-2-4k', value: '$0.04' }),
      ])
    )
  })

  test('uses the edited model name in preview row keys', () => {
    const rows = buildPreviewRows(
      {
        name: 'gpt-image-2.5-sunburst',
        price: '0.04',
        gptImage2Price1K: '0.04',
      },
      'per-request',
      '',
      '',
      '',
      EMPTY_LANE_PRICES,
      EMPTY_LANE_ENABLED,
      (key) => key
    )

    expect(rows).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          key: 'gpt-image-2.5-sunburst-1k',
          value: '$0.04',
        }),
      ])
    )
  })
})
