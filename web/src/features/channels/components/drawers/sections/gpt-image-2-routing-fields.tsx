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
import type { Control } from 'react-hook-form'
import { useTranslation } from 'react-i18next'

import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'

import type { ChannelFormValues } from '../../../lib/channel-form'

type GPTImage2RoutingFieldsProps = {
  control: Control<ChannelFormValues>
  disabled: boolean
}

export function GPTImage2RoutingFields(props: GPTImage2RoutingFieldsProps) {
  const { t } = useTranslation()
  const fields = [
    {
      name: 'gpt_image_2_model_1k' as const,
      label: t('1K upstream model'),
      description: t('Longest edge up to 1024 pixels.'),
      placeholder: 'gpt-image-2-1k',
    },
    {
      name: 'gpt_image_2_model_2k' as const,
      label: t('2K upstream model'),
      description: t('Longest edge above 1024 and up to 2048 pixels.'),
      placeholder: 'gpt-image-2-2k',
    },
    {
      name: 'gpt_image_2_model_4k' as const,
      label: t('4K / auto upstream model'),
      description: t('Auto, omitted size, or longest edge above 2048 pixels.'),
      placeholder: 'gpt-image-2-4k',
    },
  ]

  return (
    <div className='border-border/60 space-y-4 rounded-lg border p-4'>
      <div className='space-y-1'>
        <FormLabel>{t('GPT Image 2 Resolution Routing')}</FormLabel>
        <FormDescription>
          {t(
            'Route gpt-image-2 requests by the longest edge of the size field. Leave a tier empty to exclude this channel from that tier.'
          )}
        </FormDescription>
      </div>
      <div className='grid gap-4 md:grid-cols-3'>
        {fields.map((config) => (
          <FormField
            key={config.name}
            control={props.control}
            name={config.name}
            render={({ field }) => (
              <FormItem>
                <FormLabel>{config.label}</FormLabel>
                <FormControl>
                  <Input
                    {...field}
                    placeholder={config.placeholder}
                    disabled={props.disabled}
                  />
                </FormControl>
                <FormDescription>{config.description}</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
        ))}
      </div>
    </div>
  )
}
