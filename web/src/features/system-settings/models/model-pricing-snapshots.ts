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
import { splitBillingExprAndRequestRules } from '@/features/pricing/lib/billing-expr'
import { splitPluginBillingExprKey } from '@/features/pricing/lib/plugin-pricing'

import { safeJsonParse } from '../utils/json-parser'
import {
  GPT_IMAGE_2_MODEL,
  GPT_IMAGE_2_TIER_PRICE_KEYS,
} from './model-pricing-core'
import { formatPricingNumber } from './pricing-format'

export type ModelPricingSnapshotInput = {
  modelPrice: string
  modelRatio: string
  cacheRatio: string
  createCacheRatio: string
  completionRatio: string
  imageRatio: string
  audioRatio: string
  audioCompletionRatio: string
  billingMode: string
  billingExpr: string
  pluginBillingExpr?: string
}

export type ModelPricingSnapshot = {
  pluginBillingExpr?: Record<string, string>
  name: string
  price?: string
  ratio?: string
  cacheRatio?: string
  createCacheRatio?: string
  completionRatio?: string
  imageRatio?: string
  audioRatio?: string
  audioCompletionRatio?: string
  billingMode?: string
  billingExpr?: string
  requestRuleExpr?: string
  gptImage2Price1K?: string
  gptImage2Price2K?: string
  gptImage2Price4K?: string
  hasConflict: boolean
}

export type ModelRow = ModelPricingSnapshot & {
  saved?: ModelPricingSnapshot
  draft?: ModelPricingSnapshot
  isDraftChanged: boolean
  isDraftDeleted: boolean
  isDraftNew: boolean
}

export const hasPricingValue = (value?: string) =>
  value !== undefined && value !== ''

export const isBasePricingUnset = (snapshot?: ModelPricingSnapshot) =>
  !snapshot ||
  (snapshot.billingMode !== 'tiered_expr' &&
    !hasPricingValue(snapshot.price) &&
    !hasPricingValue(snapshot.ratio) &&
    !hasPricingValue(snapshot.gptImage2Price1K) &&
    !hasPricingValue(snapshot.gptImage2Price2K) &&
    !hasPricingValue(snapshot.gptImage2Price4K))

const toNumberOrNull = (value?: string) => {
  if (!hasPricingValue(value)) return null
  const num = Number(value)
  return Number.isFinite(num) ? num : null
}

const ratioToPrice = (ratio?: string, denominator?: string) => {
  const ratioNumber = toNumberOrNull(ratio)
  const denominatorNumber = denominator ? toNumberOrNull(denominator) : 2
  if (ratioNumber === null || denominatorNumber === null) return ''
  return formatPricingNumber(ratioNumber * denominatorNumber)
}

export const getModeLabel = (mode?: string) => {
  if (mode === 'per-request') return 'Per-request (deprecated)'
  if (mode === 'tiered_expr') return 'Expression'
  return 'Per-token (deprecated)'
}

export const getModeVariant = (
  mode?: string
): 'warning' | 'info' | 'success' => {
  if (mode === 'per-request') return 'warning'
  if (mode === 'tiered_expr') return 'info'
  return 'success'
}

const getExpressionSummary = (
  row: ModelPricingSnapshot,
  t: (key: string) => string
) => {
  const tierCount = (row.billingExpr?.match(/tier\(/g) || []).length
  if (tierCount > 0) {
    return `${t('Tiered pricing')} · ${tierCount} ${t('tiers')}`
  }
  return t('Expression pricing')
}

export const getPriceSummary = (
  row: ModelPricingSnapshot,
  t: (key: string) => string
) => {
  if (row.billingMode === 'tiered_expr') {
    return getExpressionSummary(row, t)
  }
  if (row.billingMode === 'per-request') {
    if (
      row.name === GPT_IMAGE_2_MODEL &&
      (row.gptImage2Price1K || row.gptImage2Price2K || row.gptImage2Price4K)
    ) {
      return [
        row.gptImage2Price1K && `1K $${row.gptImage2Price1K}`,
        row.gptImage2Price2K && `2K $${row.gptImage2Price2K}`,
        row.gptImage2Price4K && `4K $${row.gptImage2Price4K}`,
      ]
        .filter(Boolean)
        .join(' · ')
    }
    return row.price ? `$${row.price} / ${t('request')}` : t('Unset price')
  }

  const inputPrice = ratioToPrice(row.ratio)
  if (!inputPrice) return t('Unset price')

  return `${t('Input')} $${inputPrice}`
}

export const getPriceDetail = (
  row: ModelPricingSnapshot,
  t: (key: string) => string
) => {
  if (row.billingMode === 'tiered_expr') {
    return row.requestRuleExpr
      ? t('Includes request rules')
      : t('Expression based')
  }
  if (row.billingMode === 'per-request') {
    if (
      row.name === GPT_IMAGE_2_MODEL &&
      (row.gptImage2Price1K || row.gptImage2Price2K || row.gptImage2Price4K)
    ) {
      return t('Resolution-based image price')
    }
    return t('Fixed request price')
  }

  const inputPrice = ratioToPrice(row.ratio)
  if (!inputPrice) return t('No base input price')

  const details = [
    row.completionRatio &&
      `${t('Output')} $${ratioToPrice(row.completionRatio, inputPrice)}`,
    row.cacheRatio &&
      `${t('Cache')} $${ratioToPrice(row.cacheRatio, inputPrice)}`,
    row.createCacheRatio &&
      `${t('Cache write')} $${ratioToPrice(row.createCacheRatio, inputPrice)}`,
  ]
    .filter(Boolean)
    .slice(0, 2)

  return details.length > 0 ? details.join(' · ') : t('Base input price only')
}

export const buildModelSnapshots = ({
  modelPrice,
  modelRatio,
  cacheRatio,
  createCacheRatio,
  completionRatio,
  imageRatio,
  audioRatio,
  audioCompletionRatio,
  billingMode,
  billingExpr,
  pluginBillingExpr = '{}',
}: ModelPricingSnapshotInput): ModelPricingSnapshot[] => {
  const priceMap = safeJsonParse<Record<string, number>>(modelPrice, {
    fallback: {},
    context: 'model prices',
  })
  const ratioMap = safeJsonParse<Record<string, number>>(modelRatio, {
    fallback: {},
    context: 'model ratios',
  })
  const cacheMap = safeJsonParse<Record<string, number>>(cacheRatio, {
    fallback: {},
    context: 'cache ratios',
  })
  const createCacheMap = safeJsonParse<Record<string, number>>(
    createCacheRatio,
    { fallback: {}, context: 'create cache ratios' }
  )
  const completionMap = safeJsonParse<Record<string, number>>(completionRatio, {
    fallback: {},
    context: 'completion ratios',
  })
  const imageMap = safeJsonParse<Record<string, number>>(imageRatio, {
    fallback: {},
    context: 'image ratios',
  })
  const audioMap = safeJsonParse<Record<string, number>>(audioRatio, {
    fallback: {},
    context: 'audio ratios',
  })
  const audioCompletionMap = safeJsonParse<Record<string, number>>(
    audioCompletionRatio,
    { fallback: {}, context: 'audio completion ratios' }
  )
  const billingModeMap = safeJsonParse<Record<string, string>>(billingMode, {
    fallback: {},
    context: 'billing mode',
  })
  const billingExprMap = safeJsonParse<Record<string, string>>(billingExpr, {
    fallback: {},
    context: 'billing expression',
  })

  const pluginExprMap = safeJsonParse<Record<string, string>>(
    pluginBillingExpr,
    { fallback: {}, context: 'plugin billing expressions' }
  )
  const pluginExpressionsByModel = new Map<string, Record<string, string>>()
  for (const [key, expression] of Object.entries(pluginExprMap)) {
    const parts = splitPluginBillingExprKey(key)
    if (!parts) continue
    const [plugin, model] = parts
    pluginExpressionsByModel.set(model, {
      ...pluginExpressionsByModel.get(model),
      [plugin]: expression,
    })
  }
  const modelNames = new Set([
    ...pluginExpressionsByModel.keys(),
    ...Object.keys(priceMap),
    ...Object.keys(ratioMap),
    ...Object.keys(cacheMap),
    ...Object.keys(createCacheMap),
    ...Object.keys(completionMap),
    ...Object.keys(imageMap),
    ...Object.keys(audioMap),
    ...Object.keys(audioCompletionMap),
    ...Object.keys(billingModeMap),
    ...Object.keys(billingExprMap),
  ])

  const gptImage2Price1K =
    priceMap[GPT_IMAGE_2_TIER_PRICE_KEYS.price1K]?.toString() || ''
  const gptImage2Price2K =
    priceMap[GPT_IMAGE_2_TIER_PRICE_KEYS.price2K]?.toString() || ''
  const gptImage2Price4K =
    priceMap[GPT_IMAGE_2_TIER_PRICE_KEYS.price4K]?.toString() || ''
  Object.values(GPT_IMAGE_2_TIER_PRICE_KEYS).forEach((key) =>
    modelNames.delete(key)
  )
  if (gptImage2Price1K || gptImage2Price2K || gptImage2Price4K) {
    modelNames.add(GPT_IMAGE_2_MODEL)
  }
  return [...modelNames].map((name) => {
    const price = priceMap[name]?.toString() || ''
    const ratio = ratioMap[name]?.toString() || ''
    const cache = cacheMap[name]?.toString() || ''
    const createCache = createCacheMap[name]?.toString() || ''
    const completion = completionMap[name]?.toString() || ''
    const image = imageMap[name]?.toString() || ''
    const audio = audioMap[name]?.toString() || ''
    const audioCompletion = audioCompletionMap[name]?.toString() || ''
    const tierPrices =
      name === GPT_IMAGE_2_MODEL
        ? {
            gptImage2Price1K,
            gptImage2Price2K,
            gptImage2Price4K,
          }
        : {}

    const modeForModel = billingModeMap[name]
    if (modeForModel === 'tiered_expr') {
      const fullExpr = billingExprMap[name] || ''
      const { billingExpr: pureExpr, requestRuleExpr } =
        splitBillingExprAndRequestRules(fullExpr)
      return {
        name,
        pluginBillingExpr: pluginExpressionsByModel.get(name),
        billingMode: 'tiered_expr',
        billingExpr: pureExpr,
        requestRuleExpr,
        price,
        ratio,
        cacheRatio: cache,
        createCacheRatio: createCache,
        completionRatio: completion,
        imageRatio: image,
        audioRatio: audio,
        audioCompletionRatio: audioCompletion,
        ...tierPrices,
        hasConflict: false,
      }
    }

    return {
      name,
      pluginBillingExpr: pluginExpressionsByModel.get(name),
      price,
      ratio,
      cacheRatio: cache,
      createCacheRatio: createCache,
      completionRatio: completion,
      imageRatio: image,
      audioRatio: audio,
      audioCompletionRatio: audioCompletion,
      ...tierPrices,
      billingMode:
        price !== '' ||
        (name === GPT_IMAGE_2_MODEL &&
          (gptImage2Price1K !== '' ||
            gptImage2Price2K !== '' ||
            gptImage2Price4K !== ''))
          ? 'per-request'
          : 'per-token',
      hasConflict:
        price !== '' &&
        (ratio !== '' ||
          completion !== '' ||
          cache !== '' ||
          createCache !== '' ||
          image !== '' ||
          audio !== '' ||
          audioCompletion !== ''),
    }
  })
}

export const getSnapshotSignature = (snapshot?: ModelPricingSnapshot) => {
  if (!snapshot) return ''
  return JSON.stringify({
    price: snapshot.price || '',
    ratio: snapshot.ratio || '',
    cacheRatio: snapshot.cacheRatio || '',
    createCacheRatio: snapshot.createCacheRatio || '',
    completionRatio: snapshot.completionRatio || '',
    imageRatio: snapshot.imageRatio || '',
    audioRatio: snapshot.audioRatio || '',
    audioCompletionRatio: snapshot.audioCompletionRatio || '',
    billingMode: snapshot.billingMode || 'per-token',
    billingExpr: snapshot.billingExpr || '',
    requestRuleExpr: snapshot.requestRuleExpr || '',
    gptImage2Price1K: snapshot.gptImage2Price1K || '',
    gptImage2Price2K: snapshot.gptImage2Price2K || '',
    gptImage2Price4K: snapshot.gptImage2Price4K || '',
    pluginBillingExpr: Object.entries(snapshot.pluginBillingExpr ?? {}).sort(
      ([a], [b]) => a.localeCompare(b)
    ),
  })
}
