package eventstoredb

import (
    "context"
    "fmt"
    "net/http"
    "net/url"
)

func EnableByCategoryProjection(ctx context.Context, esdbUrl string) error {
    parsedUrl, err := url.Parse(esdbUrl)
    if err != nil {
        return err
    }

    req, err := http.NewRequestWithContext(
        ctx,
        http.MethodPost,
        fmt.Sprintf("http://%s:%s/projection/$by_category/command/enable", parsedUrl.Hostname(), parsedUrl.Port()),
        nil,
    )
    if err != nil {
        return fmt.Errorf("failed creating request for enabling $by_category projection: %w", err)
    }
    req.Header.Add("Accept", "application/json")
    req.Header.Add("Content-Length", "0")
    _, err = http.DefaultClient.Do(req)
    if err != nil {
        return fmt.Errorf("failed enabling $by_category projection: %w", err)
    }

    return nil
}
