# Crony

Simple cron replacement that calls URL

## Usage

Config

```json
{
  "tasks": [
    {
      "cron": "@every 1m",
      "url": "http://app/task"
    }
  ]
}

```  

Docker

 
```yaml
services:
  app:
    
  cron:
    image: 'formapro/crony:latest'
    depends_on:
      - app
    volumes:
      - './crony.json:/app/crony.json'
```

Symfony

```php
<?php
namespace App\Listener;

use Symfony\Component\EventDispatcher\EventSubscriberInterface;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\HttpKernel\Event\GetResponseEvent;
use Symfony\Component\HttpKernel\Event\PostResponseEvent;
use Symfony\Component\HttpKernel\KernelEvents;

class TaskListener implements EventSubscriberInterface
{
    public function onRequest(GetResponseEvent $event): void
    {
        if ('/task' != $event->getRequest()->getRequestUri()) {
            return;
        }

        $event->stopPropagation();
        $event->setResponse(new Response('Scheduled'));
    }

    public function onTerminate(PostResponseEvent $event): void
    {
        if ('/task' != $event->getRequest()->getRequestUri()) {
            return;
        }

        // do task
    }

    public static function getSubscribedEvents(): array
    {
        return [
            KernelEvents::REQUEST => ['onRequest', 255],
            KernelEvents::TERMINATE => 'onTerminate',
        ];
    }
}
```


## License

MIT
