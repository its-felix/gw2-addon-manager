import React from 'react';
import {Container, ContentLayout, Header} from '@cloudscape-design/components';

export function Home() {
    return (
        <ContentLayout header={<Header variant={'h1'}>Welcome to GW2 Addon Manager!</Header>}>
            <Container>Hello world</Container>
        </ContentLayout>
    );
}