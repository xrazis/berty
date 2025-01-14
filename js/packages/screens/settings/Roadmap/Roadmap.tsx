import { Layout } from '@ui-kitten/components'
import React from 'react'
import { StatusBar } from 'react-native'

import { WebViews } from '@berty/components/shared-components'
import { ScreenFC } from '@berty/navigation'
import { useThemeColor } from '@berty/store'

const RoadmapURL = 'https://guide.berty.tech/roadmap'

export const Roadmap: ScreenFC<'Settings.Roadmap'> = () => {
	const colors = useThemeColor()

	return (
		<Layout style={{ flex: 1, backgroundColor: colors['main-background'] }}>
			<StatusBar barStyle='light-content' />
			<WebViews url={RoadmapURL} />
		</Layout>
	)
}
