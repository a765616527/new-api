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

import type { Channel } from '../../types'
import {
  CHANNEL_FORM_DEFAULT_VALUES,
  transformChannelToFormDefaults,
  transformFormDataToCreatePayload,
} from '../channel-form'

describe('GPT Image 2 channel routing settings', () => {
  test('serializes configured tiers and omits blank tiers', () => {
    const result = transformFormDataToCreatePayload({
      ...CHANNEL_FORM_DEFAULT_VALUES,
      name: 'Image upstream',
      models: 'gpt-image-2',
      key: 'test-key',
      gpt_image_2_model_1k: ' vendor-image-1k ',
      gpt_image_2_model_2k: '',
      gpt_image_2_model_4k: 'vendor-image-4k',
    })

    expect(JSON.parse(result.channel.settings || '{}')).toMatchObject({
      gpt_image_2_size_models: {
        '1k': 'vendor-image-1k',
        '4k': 'vendor-image-4k',
      },
    })
  })

  test('restores all configured tiers when editing a channel', () => {
    const channel = {
      ...CHANNEL_FORM_DEFAULT_VALUES,
      id: 7,
      created_time: 0,
      test_time: 0,
      response_time: 0,
      balance: 0,
      balance_updated_time: 0,
      used_quota: 0,
      group: 'default',
      channel_info: {
        is_multi_key: false,
        multi_key_size: 0,
        multi_key_polling_index: 0,
        multi_key_mode: 'random',
      },
      settings:
        '{"gpt_image_2_size_models":{"1k":"one","2k":"two","4k":"four"}}',
    } as Channel

    const defaults = transformChannelToFormDefaults(channel)

    expect(defaults.gpt_image_2_model_1k).toBe('one')
    expect(defaults.gpt_image_2_model_2k).toBe('two')
    expect(defaults.gpt_image_2_model_4k).toBe('four')
  })

  test('serializes GPT Image 2.5 flare and sunburst tiers independently', () => {
    const result = transformFormDataToCreatePayload({
      ...CHANNEL_FORM_DEFAULT_VALUES,
      name: 'Image 2.5 upstream',
      models: 'gpt-image-2.5-flare,gpt-image-2.5-sunburst',
      key: 'test-key',
      gpt_image_2_5_flare_model_1k: 'flare-1k',
      gpt_image_2_5_sunburst_model_4k: 'sunburst-4k',
    })

    expect(JSON.parse(result.channel.settings || '{}')).toMatchObject({
      gpt_image_2_5_flare_size_models: { '1k': 'flare-1k' },
      gpt_image_2_5_sunburst_size_models: { '4k': 'sunburst-4k' },
    })
  })
})
