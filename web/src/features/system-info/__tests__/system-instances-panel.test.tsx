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
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { cleanup, render, screen, within } from '@testing-library/react'
import { afterEach, expect, it, vi } from 'vitest'

import { api } from '@/lib/api'

import { SystemInstancesPanel } from '../components/system-instances-panel'
import type { SystemInstance } from '../types'

function instance(
  nodeName: string,
  http?: { active_connections?: number }
): SystemInstance {
  return {
    node_name: nodeName,
    status: 'online',
    stale_after_seconds: 90,
    started_at: 100,
    last_seen_at: 200,
    info: { node: { name: nodeName }, http },
  }
}

function renderPanel() {
  const client = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  })
  render(
    <QueryClientProvider client={client}>
      <SystemInstancesPanel />
    </QueryClientProvider>
  )
  return client
}

afterEach(() => {
  cleanup()
  vi.restoreAllMocks()
})

it('shows each instance active request count and a dash when an older instance does not report it', async () => {
  vi.spyOn(api, 'get').mockResolvedValue({
    data: {
      success: true,
      data: [
        instance('node-busy', { active_connections: 12 }),
        instance('node-idle', { active_connections: 0 }),
        instance('node-old'),
      ],
    },
  })
  const client = renderPanel()

  const busyRow = await screen.findByRole('row', { name: /node-busy/ })
  const idleRow = screen.getByRole('row', { name: /node-idle/ })
  const oldRow = screen.getByRole('row', { name: /node-old/ })
  const headers = screen
    .getAllByRole('columnheader')
    .map((header) => header.textContent)
  const column = headers.indexOf('Active requests')
  expect(column).toBeGreaterThan(-1)

  expect(within(busyRow).getAllByRole('cell')[column]).toHaveTextContent(/^12$/)
  expect(within(idleRow).getAllByRole('cell')[column]).toHaveTextContent(/^0$/)
  expect(within(oldRow).getAllByRole('cell')[column]).toHaveTextContent(/^-$/)
  client.clear()
})
