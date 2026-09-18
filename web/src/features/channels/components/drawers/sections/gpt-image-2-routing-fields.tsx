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

type GPTImageSizeRoutingConfig = {
  model: string
  title: string
  description: string
  fields: Array<{
    name:
      | 'gpt_image_2_model_1k'
      | 'gpt_image_2_model_2k'
      | 'gpt_image_2_model_4k'
      | 'gpt_image_2_5_flare_model_1k'
      | 'gpt_image_2_5_flare_model_2k'
      | 'gpt_image_2_5_flare_model_4k'
      | 'gpt_image_2_5_sunburst_model_1k'
      | 'gpt_image_2_5_sunburst_model_2k'
      | 'gpt_image_2_5_sunburst_model_4k'
    label: string
    description: string
    placeholder: string
  }>
}

type GPTImage2RoutingFieldsProps = {
  control: Control<ChannelFormValues>
  disabled: boolean
  models?: string[]
}

export function GPTImage2RoutingFields(props: GPTImage2RoutingFieldsProps) {
  const { t } = useTranslation()
  const selected = props.models ?? ['gpt-image-2']
  const configs: GPTImageSizeRoutingConfig[] = [
    {
      model: 'gpt-image-2',
      title: t('GPT Image 2 Resolution Routing'),
      description: t(
        'Route gpt-image-2 requests by the longest edge of the size field. Leave an individual tier empty to exclude this channel from that tier. If all three tiers are empty, gpt-image-2 is passed through directly.'
      ),
      fields: [
        {
          name: 'gpt_image_2_model_1k',
          label: t('1K upstream model'),
          description: t('Longest edge up to 1536 pixels.'),
          placeholder: 'gpt-image-2-1k',
        },
        {
          name: 'gpt_image_2_model_2k',
          label: t('2K upstream model'),
          description: t('Longest edge above 1536 and up to 2048 pixels.'),
          placeholder: 'gpt-image-2-2k',
        },
        {
          name: 'gpt_image_2_model_4k',
          label: t('4K / auto upstream model'),
          description: t(
            'Auto, omitted size, or longest edge above 2048 pixels.'
          ),
          placeholder: 'gpt-image-2-4k',
        },
      ],
    },
    {
      model: 'gpt-image-2.5-flare',
      title: t('GPT Image 2.5 Flare Resolution Routing'),
      description: t(
        'Route gpt-image-2.5-flare requests by the longest edge of the size field. Leave an individual tier empty to exclude this channel from that tier. If all three tiers are empty, gpt-image-2.5-flare is passed through directly.'
      ),
      fields: [
        {
          name: 'gpt_image_2_5_flare_model_1k',
          label: t('1K upstream model'),
          description: t('Longest edge up to 1536 pixels.'),
          placeholder: 'gpt-image-2.5-flare-1k',
        },
        {
          name: 'gpt_image_2_5_flare_model_2k',
          label: t('2K upstream model'),
          description: t('Longest edge above 1536 and up to 2048 pixels.'),
          placeholder: 'gpt-image-2.5-flare-2k',
        },
        {
          name: 'gpt_image_2_5_flare_model_4k',
          label: t('4K / auto upstream model'),
          description: t(
            'Auto, omitted size, or longest edge above 2048 pixels.'
          ),
          placeholder: 'gpt-image-2.5-flare-4k',
        },
      ],
    },
    {
      model: 'gpt-image-2.5-sunburst',
      title: t('GPT Image 2.5 Sunburst Resolution Routing'),
      description: t(
        'Route gpt-image-2.5-sunburst requests by the longest edge of the size field. Leave an individual tier empty to exclude this channel from that tier. If all three tiers are empty, gpt-image-2.5-sunburst is passed through directly.'
      ),
      fields: [
        {
          name: 'gpt_image_2_5_sunburst_model_1k',
          label: t('1K upstream model'),
          description: t('Longest edge up to 1536 pixels.'),
          placeholder: 'gpt-image-2.5-sunburst-1k',
        },
        {
          name: 'gpt_image_2_5_sunburst_model_2k',
          label: t('2K upstream model'),
          description: t('Longest edge above 1536 and up to 2048 pixels.'),
          placeholder: 'gpt-image-2.5-sunburst-2k',
        },
        {
          name: 'gpt_image_2_5_sunburst_model_4k',
          label: t('4K / auto upstream model'),
          description: t(
            'Auto, omitted size, or longest edge above 2048 pixels.'
          ),
          placeholder: 'gpt-image-2.5-sunburst-4k',
        },
      ],
    },
  ]

  const visible = configs.filter((config) => selected.includes(config.model))
  if (visible.length === 0) return null

  return (
    <div className='space-y-4'>
      {visible.map((config) => (
        <div
          key={config.model}
          className='border-border/60 space-y-4 rounded-lg border p-4'
        >
          <div className='space-y-1'>
            <FormLabel>{config.title}</FormLabel>
            <FormDescription>{config.description}</FormDescription>
          </div>
          <div className='grid gap-4 md:grid-cols-3'>
            {config.fields.map((fieldConfig) => (
              <FormField
                key={fieldConfig.name}
                control={props.control}
                name={fieldConfig.name}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{fieldConfig.label}</FormLabel>
                    <FormControl>
                      <Input
                        {...field}
                        placeholder={fieldConfig.placeholder}
                        disabled={props.disabled}
                      />
                    </FormControl>
                    <FormDescription>{fieldConfig.description}</FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            ))}
          </div>
        </div>
      ))}
    </div>
  )
}
