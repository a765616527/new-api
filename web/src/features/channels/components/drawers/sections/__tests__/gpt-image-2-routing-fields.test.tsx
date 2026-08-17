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
import { useForm } from 'react-hook-form'
import { describe, expect, test } from 'vitest'

import { Form } from '@/components/ui/form'

import {
  CHANNEL_FORM_DEFAULT_VALUES,
  type ChannelFormValues,
} from '../../../../lib/channel-form'
import { GPTImage2RoutingFields } from '../gpt-image-2-routing-fields'

function Harness(props: { disabled?: boolean }) {
  const form = useForm<ChannelFormValues>({
    defaultValues: CHANNEL_FORM_DEFAULT_VALUES,
  })

  return (
    <Form {...form}>
      <GPTImage2RoutingFields
        control={form.control}
        disabled={props.disabled === true}
      />
    </Form>
  )
}

describe('GPT Image 2 resolution routing fields', () => {
  test('shows one upstream model input for every supported tier', () => {
    render(<Harness />)

    expect(
      screen.getByRole('textbox', { name: '1K upstream model' })
    ).toBeVisible()
    expect(
      screen.getByRole('textbox', { name: '2K upstream model' })
    ).toBeVisible()
    expect(
      screen.getByRole('textbox', { name: '4K / auto upstream model' })
    ).toBeVisible()
  })

  test('disables every input while the channel is being saved', () => {
    render(<Harness disabled />)

    for (const input of screen.getAllByRole('textbox')) {
      expect(input).toBeDisabled()
    }
  })
})
