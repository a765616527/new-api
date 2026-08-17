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
import { render, screen } from '@testing-library/react'
import i18next from 'i18next'
import { beforeAll, describe, expect, test } from 'vitest'

import type { UsageLog } from '../../data/schema'
import { DetailsDialog } from '../dialogs/details-dialog'

describe('GPT Image 2 request log details', () => {
  beforeAll(() => {
    i18next.addResourceBundle('en', 'translation', {
      'Image Request': 'Image Request',
      'Requested Size': 'Requested Size',
      'Resolution Tier': 'Resolution Tier',
      'Requested Quality': 'Requested Quality',
      'Image Count': 'Image Count',
      'Unit Price': 'Unit Price',
      'per image': 'per image',
    })
  })

  test('shows structured size, tier, quality, count, and unit price', () => {
    const log: UsageLog = {
      id: 1,
      user_id: 1,
      created_at: 1,
      type: 2,
      content: '',
      username: 'user',
      token_name: 'token',
      model_name: 'gpt-image-2',
      quota: 1000,
      prompt_tokens: 0,
      completion_tokens: 0,
      use_time: 1,
      is_stream: false,
      channel: 1,
      channel_name: 'image',
      token_id: 1,
      group: 'default',
      ip: '',
      other: JSON.stringify({
        image_size: '1536x1024',
        image_size_tier: '2K',
        image_quality: 'high',
        image_count: 2,
        image_unit_price: 0.04,
      }),
      request_id: 'request-id',
      upstream_request_id: '',
    }

    render(
      <DetailsDialog log={log} isAdmin open onOpenChange={() => undefined} />
    )

    expect(screen.getByText('Image Request')).toBeInTheDocument()
    expect(screen.getByText('1536x1024')).toBeInTheDocument()
    expect(screen.getByText('2K')).toBeInTheDocument()
    expect(screen.getByText('high')).toBeInTheDocument()
    expect(screen.getByText(/0\.04.*per image/)).toBeInTheDocument()
  })
})
